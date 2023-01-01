package main

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-22 10:19:45
 * @Description:
 */

import (
	"chatshock/configs"
	"chatshock/middlewares"
	"chatshock/websockets"
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceQuote:      true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	var configVersion string
	flag.StringVar(&configVersion, "c", "dev", "请输入配置版本(dev,product,)")
	flag.Parse()
	var r = gin.New()
	r.Use(middlewares.AccessLog())
	r.Use(gin.Recovery())
	InitConfig(configVersion)
	InitDatabase()
	InitRedis()
	InitMinioClient()
	InitRouters(r)

	err := websockets.BroadCaster.ReLinkChatRooms()
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}
	runHost := configs.Conf.AppConfig.AppHost + ":" + configs.Conf.AppConfig.AppPort
	err = r.Run(runHost)
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}
}
