package main

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 14:16:34
 * @Description:
 */

import (
	"chatshock/configs"
	"chatshock/middlewares"
	"flag"
	"os"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	var configVersion string
	flag.StringVar(&configVersion, "c", "dev", "请输入配置版本(dev,product,)")
	flag.Parse()
	var r *gin.Engine
	r = gin.New()
	r.Use(middlewares.AccessLog())
	r.Use(middlewares.Recovery())
	InitConfig(configVersion)
	InitDatabase()
	InitRedis()
	InitMinioClient()
	InitSmsClient()
	InitRouters(r)
	runHost := configs.Conf.AppConfig.AppHost + ":" + configs.Conf.AppConfig.AppPort
	err := r.Run(runHost)
	if err != nil {
		log.Fatal(err.Error())
	}
}
