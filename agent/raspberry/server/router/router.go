package router

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/start", service.Start)
	router.POST("/run", service.Run)
	router.POST("/gatewayIP", service.SetGatewayIP)
	router.POST("/stop", service.Stop)
	router.POST("/service", service.UpdateService)
	router.DELETE("/service", service.DeleteService)
	router.GET("/service", service.GetService)
	return router
}
