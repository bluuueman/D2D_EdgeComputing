package router

import (
	"raspberry/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/service", service.PostServiceInfo)
	router.POST("/server", service.PostServerInfo)
	router.POST("/job", service.PostJob)
	router.DELETE("/job", service.DeleteJob)
	router.POST("/echoTime", service.EchoTime)
	router.POST("/postTest", service.PostTest)
	return router
}
