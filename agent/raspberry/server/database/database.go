package database

import (
	"errors"
	"log"
	"net"
)

var run bool
var gateway_ip string
var local_ip string
var ports map[string]string
var cmds map[string]string

func Init() {
	run = false
	gateway_ip = ""
	setLocalIP()
	ports = make(map[string]string)
	cmds = make(map[string]string)
}

//获取ip
func getExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

//获取ip
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func SetGatewayIP(ip string) {
	gateway_ip = ip
}

func GetGatewayIP() string {
	return gateway_ip
}

func GetLocalIP() string {
	return local_ip
}

func setLocalIP() {
	ip, err := getExternalIP()
	if err != nil {
		log.Println(err)
	}
	local_ip = ip.String()
}

func SetRun() {
	run = true
}

func SetStop() {
	run = false
}

func IsRun() bool {
	return run
}

func SetService(service string, port string, cmd string) {
	ports[service] = port
	cmds[service] = cmd
}

func DeleteService(service string) {
	_, exist := ports[service]
	if exist {
		delete(ports, service)
		delete(cmds, service)
	}
}

func GetService() map[string]string {
	return ports
}
