package database

import (
	"fmt"
	"raspberry/stream"
	"raspberry/utility"
	"sync"
	"time"
)

var cmd chan [2]string

type JobInfo struct {
	//key is service name
	//channel to change job selecting thread status
	channel map[string](chan int)
	//0 pending 1 selecting 2 running 3 stop
	status map[string]int
	//ip of the server that service offloading to
	ip   map[string][]string
	lock sync.RWMutex
}

var ji JobInfo

//Init jobinfo data struct
func InitJobInfo() {
	channel := make(map[string](chan int))
	status := make(map[string]int)
	ip := make(map[string][]string)
	ji.channel = channel
	ji.status = status
	ji.ip = ip
}

//Add job you want the gateway to offload
func AddJob(service string) int {
	ji.lock.Lock()
	defer ji.lock.Unlock()
	_, exist := ji.status[service]
	if !exist {
		//after add it to jobinfo, start thread then change job status to 1(selecting)
		ch := make(chan int)
		ji.channel[service] = ch
		ji.status[service] = 1
		return 1
	}
	return 0
}

//After a needed service has been registered, change job status to 1(selecting)
func WakeJob(service string) {
	ji.lock.Lock()
	defer ji.lock.Unlock()
	status, exist := ji.status[service]
	if exist && status == 0 {
		ji.channel[service] <- 1
		ji.status[service] = 1
	}

}

func PrintJobStatus() {
	ji.lock.RLock()
	for service, status := range ji.status {
		fmt.Println(service, " ", status)
	}
}

//Check all available server heartbeat
func KeepAliveAll() {
	for {
		fmt.Println("check server")
		//if timestamp outdate, delete server info and service info
		CheckServerStatus()
		PrintServerInfo()
		PrintServerlist()
		fmt.Println()
		time.Sleep(time.Duration(10) * time.Second)
	}
}

//Check runing server heartbeat
func KeepAliveService() {
	for {
		ji.lock.Lock()
		fmt.Println("check Joblist")
		checked := make(map[string]int)
		for service, status := range ji.status {
			fmt.Println(service, " ", status)
			//if the job has been stoped, wake the thread
			if status == 3 {
				ji.channel[service] <- 0
			} else if status == 2 { //if the job is runing check all server heartbeat
				si.lock.RLock()
				now := time.Now().Unix()
				for _, ip := range ji.ip[service] {
					_, exist := checked[ip]
					if exist {
						if checked[ip] == 1 {
							break
						} else {
							si.priority[ip] = 0
							ji.status[service] = 1
							ji.channel[service] <- 1
							break
						}
					}
					if now-si.status[ip] > 4 {
						si.priority[ip] = -1
						ji.status[service] = 1
						ji.channel[service] <- 1
						//fmt.Println(service, " seleting")
					} else {
						checked[ip] = 1
					}
				}
				si.lock.RUnlock()
			}
		}
		fmt.Println()
		ji.lock.Unlock()
		time.Sleep(time.Duration(4) * time.Second)
	}
}

//Thread to monitor job status, need to be waked after job status change to 1 or 3
func SelectServer(service string, ch chan int) {
	var run *bool
	run = nil
	for {
		ji.lock.Lock()
		if ji.status[service] == 3 {
			if run != nil {
				*run = false
			}
			//if job is stoped, remove it data
			delete(ji.channel, service)
			delete(ji.status, service)
			delete(ji.ip, service)
			//fmt.Println(service, " delete")
			//fmt.Println()
			ji.lock.Unlock()
			return
		} else {
			if run != nil {
				*run = false
			}
			//try to find available server
			result := GetService(service)
			ji.ip[service] = []string{}
			//PrintServerlist()
			//if find one, change status to 2(running)
			if result[0] != "" {
				run = new(bool)
				*run = true
				//fmt.Println(service, " runing")
				ji.status[service] = 2
				ji.ip[service] = append(ji.ip[service], result[0])
				go utility.NoticeServer(ji.ip[service][0], service)
				url := "http://192.168.0.168:5002/detect"
				go stream.Streamer(url, 300*time.Millisecond, run)
				//updat timestamp here
			} else {
				//if not, change status to pending
				//fmt.Println(service, " pending")
				ji.status[service] = 0
			}
		}
		fmt.Println()
		ji.lock.Unlock()
		<-ch
	}

}

//Start job selecting thread
func StartJob(service string) {
	ji.lock.Lock()
	defer ji.lock.Unlock()
	_, exist := ji.status[service]
	if !exist {
		ji.status[service] = 1
		ch := make(chan int, 2)
		ji.channel[service] = ch
		//start selecting thread to monitor job status
		go SelectServer(service, ch)
	}

}

//Delete job gateway want to offload
func StopJob(service string) {
	ji.lock.Lock()
	defer ji.lock.Unlock()
	_, exist := ji.status[service]
	if exist {
		ji.status[service] = 3
	}
}

//not used
func MainRuntime() {
	run := "run"
	stop := "stop"
	for {
		cmd := <-cmd
		opt := cmd[0]
		target := cmd[1]
		switch opt {
		case run:
			ji.lock.Lock()
			_, exist := ji.status[target]
			if !exist {
				ji.status[target] = 1
				ch := make(chan int)
				ji.channel[target] = ch
				go SelectServer(target, ch)
			}
			ji.lock.Unlock()
		case stop:
			ji.lock.Lock()
			_, exist := ji.status[target]
			if exist {
				ji.status[target] = 3
			}
			ji.lock.Unlock()
		default:
			fmt.Println("wrong cmd")
		}
	}

}
