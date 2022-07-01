package service

import (
	"fmt"
	"net/http"
	"server/database"
	"server/utility"
	"time"

	"github.com/gin-gonic/gin"
)

func Start(c *gin.Context) {
	type msg struct {
		Service string `json:"service"`
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
	utility.StartService(jsondata.Service)
	c.JSON(http.StatusOK, gin.H{
		"message": "Message receive",
	})

}

func UpdateService(c *gin.Context) {
	type msg struct {
		Service string `json:"service"`
		Port    string `json:"port"`
		Cmd     string `json:"cmd"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	database.SetService(jsondata.Service, jsondata.Port, jsondata.Cmd)
	c.JSON(http.StatusOK, gin.H{
		"message": "Service Update",
	})
}

func DeleteService(c *gin.Context) {
	type msg struct {
		Service string `json:"service"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	database.DeleteService(jsondata.Service)
	c.JSON(http.StatusOK, gin.H{
		"message": "Service Delete",
	})
}

func GetService(c *gin.Context) {
	services := database.GetService()
	c.JSON(http.StatusOK, gin.H{
		"data": services,
	})
}

func SetGatewayIP(c *gin.Context) {
	type msg struct {
		IP string `json:"gatewayIP"`
	}
	jsondata := msg{}
	bindErr := c.BindJSON(&jsondata)
	if utility.IsErr(bindErr, "BindJSON Failed!") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server JSON bind failed",
		})
		return
	}
	database.SetGatewayIP(jsondata.IP)
	c.JSON(http.StatusOK, gin.H{
		"message": "Gateway IP changed",
	})
}

func Run(c *gin.Context) {
	ip := database.GetGatewayIP()
	if ip == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gateway IP has not been seted yet",
		})
		return
	}
	if database.IsRun() {
		c.JSON(http.StatusOK, gin.H{
			"message": "Already run",
		})
		return
	}
	database.SetRun()
	utility.RegisterService()
	go utility.SendHeartBeat()
	c.JSON(http.StatusOK, gin.H{
		"message": "Runing",
	})
}

func Stop(c *gin.Context) {
	if database.IsRun() {
		database.SetStop()
		c.JSON(http.StatusOK, gin.H{
			"message": "Stoped",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Notihing runing now",
		})
	}
}
