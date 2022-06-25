package main

import (
	"raspberry/database"
	"raspberry/router"
	"raspberry/stream"
)

func main() {
	router := router.SetupRouter()
	database.InitAll()
	stream.InitStream()
	go database.KeepAliveAll()
	go database.KeepAliveService()
	_ = router.Run(":8080")
	stream.StopStream()
}
