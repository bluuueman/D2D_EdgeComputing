package main

import (
	"server/database"
	"server/router"
)

func main() {
	router := router.SetupRouter()
	services := make(map[string]string)
	services["test1"] = "8080"
	services["test2"] = "8000"
	database.Init()
	//data := utility.GetService(services)
	//desIp := "192.168.0.164"
	//srcIp := "192.168.0.168"
	//utility.RegisterService(desIp, srcIp, data, 3)
	//go utility.SendHeartBeat(desIp, srcIp)
	_ = router.Run(":8000")

}
