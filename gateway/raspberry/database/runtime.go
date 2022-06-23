package database

import (
	"fmt"
	"sync"
	"time"
)

var cmd chan [2]string

type JobInfo struct {
	channel map[string](chan int)
	status  map[string]int
	ip      map[string]string
	lock    sync.RWMutex
}

var ji JobInfo

func InitJobInfo() {
	channel := make(map[string](chan int))
	status := make(map[string]int)
	ip := make(map[string]string)
	ji.channel = channel
	ji.status = status
	ji.ip = ip
}

func AddJob(service string) int {
	ji.lock.Lock()
	defer ji.lock.Unlock()
	_, exist := ji.status[service]
	if !exist {
		ch := make(chan int)
		ji.channel[service] = ch
		ji.status[service] = 1
		return 1
	}
	return 0
}

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

func KeepAliveAll() {
	for {
		fmt.Println("check server")
		CheckServerStatus()
		//PrintServerInfo()
		//PrintServerlist()
		//fmt.Println()
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func KeepAliveService() {
	for {
		ji.lock.Lock()
		fmt.Println("check Joblist")
		for service, status := range ji.status {
			//fmt.Println(service, " ", status)
			if status == 3 {
				ji.channel[service] <- 0
			} else if status == 2 {
				si.lock.RLock()
				now := time.Now().Unix()
				ip := ji.ip[service]
				if now-si.status[ip] > 4 {
					si.priority[ip] = 0
					ji.status[service] = 1
					ji.channel[service] <- 1
					//fmt.Println(service, " seleting")
				}
				si.lock.RUnlock()
			}
		}
		//fmt.Println()
		ji.lock.Unlock()
		time.Sleep(time.Duration(4) * time.Second)
	}
}

func SelectServer(service string, ch chan int) {
	for {
		ji.lock.Lock()
		if ji.status[service] == 3 {
			delete(ji.channel, service)
			delete(ji.status, service)
			delete(ji.ip, service)
			//fmt.Println(service, " delete")
			//fmt.Println()
			ji.lock.Unlock()
			return
		} else {
			result := GetService(service)
			//PrintServerlist()
			if result[0] != "" {
				//fmt.Println(service, " runing")
				ji.status[service] = 2
				ji.ip[service] = result[0]
				//updat timestamp here
			} else {
				//fmt.Println(service, " pending")
				ji.status[service] = 0
			}
		}
		fmt.Println()
		ji.lock.Unlock()
		<-ch
	}

}

func StartJob(service string) {
	ji.lock.Lock()
	defer ji.lock.Unlock()
	_, exist := ji.status[service]
	if !exist {
		ji.status[service] = 1
		ch := make(chan int, 2)
		ji.channel[service] = ch
		go SelectServer(service, ch)
	}

}

func StopJob(service string) {
	ji.lock.Lock()
	defer ji.lock.Unlock()
	_, exist := ji.status[service]
	if exist {
		ji.status[service] = 3
	}
}

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
