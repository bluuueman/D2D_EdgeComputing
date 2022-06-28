package utility

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func IsErr(err error, msg string) bool {
	if err != nil {
		log.Println("ERROR: "+msg+"\n", err)
		return true
	}
	return false
}
func HttpSend(url string, data string) {
	var jsonStr = []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	IsErr(err, "Http Request Generate Failed")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if IsErr(err, "Http Send Failed") {
		return
	}
	defer resp.Body.Close()
	fmt.Println("status", resp.Status)
}
func RegisterService(desIp string, srcIp string, service string, priority int) {
	time.Sleep(time.Duration(3) * time.Second)
	url := "http://" + desIp + ":8080/service"
	//json序列化
	tmp := strconv.Itoa(priority)
	data := "{\"priority\":" + tmp +
		",\"ip\":\"" + srcIp +
		"\",\"data\":" + service +
		"}"
	HttpSend(url, data)
}
func HeartBeat(desIp string, srcIp string, priority int) {
	url := "http://" + desIp + ":8080/server"
	tmp := strconv.Itoa(priority)
	data := "{\"priority\":" + tmp +
		",\"ip\":\"" + srcIp +
		"\"}"
	HttpSend(url, data)
}
func SendHeartBeat(desIp string, srcIp string) {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		HeartBeat(desIp, srcIp, 3)
	}
}
func GetService(services map[string]string) string {
	i := 1
	data := "{"
	for service, port := range services {
		tmp := "{\"service\":" + "\"" + service + "\"" +
			",\"port\":" + "\"" + port + "\"" +
			"}"
		if i == 1 {
			data += "\"" + strconv.Itoa((i)) + "\":" + tmp
		} else {
			data += ",\"" + strconv.Itoa((i)) + "\":" + tmp
		}
		i++
	}
	data += "}"
	return data
}

func WriteMessage(conn net.Conn, msg string) (int, error) {
	fmt.Println("Sending ", msg)
	var buf bytes.Buffer
	buf.WriteString(msg)
	return conn.Write(buf.Bytes())
}

func StartService(url string, port string) {
	/*
		conn, err := net.Dial("tcp", "127.0.0.1:"+port)
		IsErr(err, "Connect Service Failed")
		defer conn.Close()
		fmt.Println("Connect to Service")
		_, w_err := WriteMessage(conn, url)
		IsErr(w_err, "Write Failed")
	*/
	fmt.Println("Start service....." + url + ":" + port)

}
