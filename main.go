package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"chatshock/configs"

)

func main() {
	r := gin.Default()
	InitConfig()
	InitDatabase()
	InitRedis()
	InitRouters(r)
	runHost := configs.Conf.AppConfig.AppHost + ":" + configs.Conf.AppConfig.AppPort
	err := r.Run(runHost)
	if err != nil {
		log.Fatal(err.Error())
	}
}