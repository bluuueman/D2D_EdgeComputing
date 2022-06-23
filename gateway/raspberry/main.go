package main

import (
	"raspberry/database"
	"raspberry/router"
)

func main() {
	router := router.SetupRouter()
	database.InitAll()
	go database.KeepAliveAll()
	go database.KeepAliveService()
	_ = router.Run(":8080")
}
