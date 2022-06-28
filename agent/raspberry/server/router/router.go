package router

import (
	"server/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/postTest", service.PostTest)
	return router
}
