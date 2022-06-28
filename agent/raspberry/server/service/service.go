package service

import (
	"fmt"
	"net/http"
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
