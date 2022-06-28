package service

import (
	"fmt"
	"net/http"
	"server/utility"
	"time"

	"github.com/gin-gonic/gin"
)

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
	utility.StartService("192.168.0.168:8080", "1234")
	c.JSON(http.StatusOK, gin.H{
		"message": "Message receive",
	})

}
