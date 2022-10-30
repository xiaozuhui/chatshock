package main

import (
	"chatshock/configs"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	var configVersion string
	flag.StringVar(&configVersion, "c", "dev", "请输入配置版本(dev,product,)")
	flag.Parse()
	r := gin.Default()
	InitConfig()
	InitDatabase()
	InitRedis()
	InitMinioClient()
	InitRouters(r)
	runHost := configs.Conf.AppConfig.AppHost + ":" + configs.Conf.AppConfig.AppPort
	err := r.Run(runHost)
	if err != nil {
		log.Fatal(err.Error())
	}
}
