package service

import (
	"fmt"
	"net/http"
	"raspberry/database"
	"raspberry/utility"
	"time"

	"github.com/gin-gonic/gin"
)

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
                "type": jsondata.Type,
		"message": time.Now().UnixMilli(),
	})
	return
}

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
