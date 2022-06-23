package database

import (
	"fmt"
	"raspberry/utility"
	"sync"
	"time"
)

type Server struct {
	len     int
	service map[string]string
	lock    sync.RWMutex
}

type ServerList struct {
	len    int
	server map[string]*Server
	lock   sync.RWMutex
}

type ServerInfo struct {
	len      int
	status   map[string]int64
	priority map[string]int
	lock     sync.RWMutex
}

var sl ServerList
var si ServerInfo

func InitServer() *Server {
	s := new(Server)
	service := make(map[string]string)
	s.service = service
	s.len = 0
	return s
}

func UpdateService(server *Server, service string, port string) {
	server.lock.Lock()
	defer server.lock.Unlock()
	_, exist := server.service[service]
	server.service[service] = port
	if !exist {
		server.len += 1
	}
}

func DeletService(server *Server, service string) {
	server.lock.Lock()
	defer server.lock.Unlock()
	_, exist := server.service[service]
	if exist {
		delete(server.service, service)
		server.len -= 1
	}
}

func SearchService(server *Server, service string) string {
	server.lock.RLock()
	defer server.lock.RUnlock()
	port, exist := server.service[service]
	if exist {
		return port
	} else {
		return ""
	}
}

/***********************************************************************************************************






***********************************************************************************************************/

func InitServerList() {
	server := make(map[string]*Server)
	sl.server = server
	sl.len = 0
}

func AddServer(ip string, pri int) *Server {
	UpdataServerInfo(ip, pri)
	sl.lock.Lock()
	defer sl.lock.Unlock()
	_, exist := sl.server[ip]
	if !exist {
		server := InitServer()
		sl.server[ip] = server
		sl.len += 1
	}
	return sl.server[ip]

}

func DeletServer(ip string) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	_, exist := sl.server[ip]
	if exist {
		delete(sl.server, ip)
		sl.len -= 1
	}
}

func GetServer(ip string) *Server {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	server, exist := sl.server[ip]
	if exist {
		return server
	} else {
		return nil
	}
}

func PrintServerlist() {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	fmt.Println(sl.len)
	for ip, server := range sl.server {
		server.lock.RLock()
		fmt.Println(ip, server.len)
		for service, port := range server.service {
			fmt.Println("Service Name: ", service, " Port: ", port)
		}
		server.lock.RUnlock()
	}
}

/***********************************************************************************************************






***********************************************************************************************************/

func InitServerInfo() {
	status := make(map[string]int64)
	priority := make(map[string]int)
	si.len = 0
	si.priority = priority
	si.status = status
}

func UpdataServerInfo(ip string, pri int) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.status[ip] = time.Now().Unix()
	si.priority[ip] = pri
	si.len = len(si.status)
}

func PrintServerInfo() {
	si.lock.RLock()
	defer si.lock.RUnlock()
	count := 0
	for ip, timestamp := range si.status {
		fmt.Println("Ip: ", ip, " Timestamp: ", timestamp)
		fmt.Println("Priority: ", si.priority[ip])
		count += 1
	}
	if count == 0 {
		fmt.Println("Empty Server Status")
	}
}

/***********************************************************************************************************






***********************************************************************************************************/

func CheckServerStatus() {
	si.lock.Lock()
	defer si.lock.Unlock()
	if si.len > 0 {
		now := time.Now().Unix()
		for ip, timestamp := range si.status {
			if now-timestamp > 10 {
				delete(si.status, ip)
				si.len -= 1
				DeletServer(ip)
			}
		}
	}
}

func CheckServiceServer(ip string) {

}

func InitAll() {
	InitServerList()
	InitServerInfo()
	InitJobInfo()
}

func GetService(service string) [2]string {
	si.lock.RLock()
	serverlist := utility.Rank(si.priority)
	fmt.Println(serverlist)
	si.lock.RUnlock()
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	for _, ip := range serverlist {
		server, exist := sl.server[ip]
		if exist {
			port := SearchService(server, service)
			if port != "" {
				result := [2]string{ip, port}
				return result
			}
		}
	}
	result := [2]string{"", ""}
	return result
}
