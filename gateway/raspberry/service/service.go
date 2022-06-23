package service

import (
	"fmt"
	"net/http"
	"raspberry/database"
	"raspberry/utility"
	"time"

	"github.com/gin-gonic/gin"
)

/*Update server's service info
* URL             ip:port/service
* Method          POST
* Content-Type    application/json
* Body
* {
*     "ip":"192.168.0.1",
*     "priority":4,
*     "data":{
*         "1":{
*             "service":"s1",
*             "port":"8080"
*         },
*         "2":{
*             "service":"s2",
*              "prot":"8088"
*         }
*     }
* }
*
 */
func PostServiceInfo(c *gin.Context) {
	type Service struct {
		ServiceName string `json:"service"`
		Port        string `json:"port"`
	}
	type msg struct {
		Data     map[string]Service `json:"data"`
		IP       string             `json:"ip"`
		Priority int                `json:"priority"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	server := database.GetServer(jsondata.IP)
	if server == nil {
		server = database.AddServer(jsondata.IP, jsondata.Priority)
	}
	for _, item := range jsondata.Data {
		database.UpdateService(server, item.ServiceName, item.Port)
		database.WakeJob(item.ServiceName)
	}
	fmt.Println("Recive a service register")
	fmt.Println(time.Now().UnixMilli())
	go utility.NoticeServer(jsondata.IP, "test", "test")
	c.JSON(http.StatusOK, gin.H{
		"message": "Service Info Recv",
	})
}

/*Update server status as heartbeat
* URL             ip:port/server
* Method          POST
* Content-Type    application/json
* Body
* {
*     "ip":"192.168.0.1",
*     "priority":4
* }
 */
func PostServerInfo(c *gin.Context) {
	type msg struct {
		IP       string `json:"ip"`
		Priority int    `json:"priority"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	database.UpdataServerInfo(jsondata.IP, jsondata.Priority)

	c.JSON(http.StatusOK, gin.H{
		"message": "Server Info Recv",
	})
}

/*Add service you want gateway offload
* URL             ip:port/job
* Method          POST
* Content-Type    application/json
* Body
* {
*     "service":"service name"
* }
 */
func PostJob(c *gin.Context) {
	type msg struct {
		ServiceName string `json:"service"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	database.StartJob(jsondata.ServiceName)
	c.JSON(http.StatusOK, gin.H{
		"message": "Job add",
	})
}

/*Delete service you dont want gateway offload
* URL             ip:port/server
* Method          Delete
* Content-Type    application/json
* Body
* {
*     "service":"service name"
* }
 */
func DeleteJob(c *gin.Context) {
	type msg struct {
		ServiceName string `json:"service"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	database.StopJob(jsondata.ServiceName)
	c.JSON(http.StatusOK, gin.H{
		"message": "Job stop",
	})
}

/*For test only :echo time
* URL             ip:port/echotime
* Method          POST
* Content-Type    application/json
* Body
* {
*     "type":"type name"
* }
 */
func EchoTime(c *gin.Context) {
	type msg struct {
		Type string `json:"type"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	fmt.Println("receive message")
	fmt.Println(time.Now().UnixMilli())
	c.JSON(http.StatusOK, gin.H{
		"type":    jsondata.Type,
		"message": time.Now().UnixMilli(),
	})
	return
}

/*For test only :receive offload request
* URL             ip:port/echotime
* Method          POST
* Content-Type    application/json
* Body
* {
*     "service":"service name",
*     "port":"8080"
* }
 */
func PostTest(c *gin.Context) {
	type msg struct {
		Service string `json:"service"`
		Port    string `json:"port"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	fmt.Println("receive message")
	fmt.Println(time.Now().UnixMilli())
	c.JSON(http.StatusOK, gin.H{
		"message": "Message receive",
	})
	return
}
