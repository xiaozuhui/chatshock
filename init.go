package main

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-09 13:44:46
 * @Description: 初始化各种配置
 */

import (
	"chatshock/configs"
	"chatshock/models"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDatabase 初始化数据库
/**
 * @description: 初始化数据库
 * @return {*}
 * @author: xiaozuhui
 */
func InitDatabase() {
	log.Info("初始化数据库开始...")
	var err error
	dbConfig := configs.Conf.DBConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBName, dbConfig.DBPort, dbConfig.SSLMode)
	configs.DBEngine, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(errors.WithStack(err))
	}
	if configs.DBEngine == nil {
		panic(errors.WithStack(errors.New("DB 初始化失败.")))
	}
	err = models.InitModel()
	if err != nil {
		log.Error(err.Error())
		panic(errors.WithStack(err))
	}
	log.Info("初始化数据库结束.")
}

// InitConfig 初始化配置，获取并解析configs中的配置文件
/**
 * @description: 初始化配置，获取并解析configs中的配置文件
 * @param {string} configVersion
 * @return {*}
 * @author: xiaozuhui
 */
func InitConfig(configVersion string) *configs.Config {
	log.Info("初始化Config配置...")
	configs.Conf = &configs.Config{}
	configs.BaseDir, _ = os.Getwd()
	viper.SetConfigName(configVersion)
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(configs.BaseDir, "configs", configVersion+".yaml"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(errors.WithStack(err))
	}
	configs.Conf.Parse(viper.GetViper())
	log.Info("初始化Config配置完毕.")
	log.Infof("获取的配置为：%s", configs.Conf.String())
	return configs.Conf
}

// InitRedis 初始化redis
/**
 * @description: 初始化redis
 * @return {*}
 * @author: xiaozuhui
 */
func InitRedis() *redis.Client {
	log.Info("初始化Redis客户端...")
	redisConfig := configs.Conf.RedisConfig
	configs.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.RedisHost, redisConfig.RedisPort),
		Password: "",
	})
	result := configs.RedisClient.Ping(context.Background())
	if result.Val() != "PONG" {
		// 连接有问题
		fmt.Println("redis链接不上")
		return nil
	}
	log.Info("Redis客户端初始化完毕.")
	return configs.RedisClient
}

// InitMinioClient 初始化minio，并获取其client
/**
 * @description: 初始化minio，并获取其client
 * @return {*}
 * @author: xiaozuhui
 */
func InitMinioClient() *minio.Client {
	log.Info("初始化minio客户端...")
	var err error
	configs.MinioClient, err = minio.New(configs.Conf.MinioConfig.EndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(configs.Conf.MinioConfig.AccessKeyID,
			configs.Conf.MinioConfig.SecretAccessKey, ""),
		Secure: configs.Conf.MinioConfig.UseSSL,
	})
	if err != nil {
		log.Fatal(errors.WithStack(err))
		panic(errors.WithStack(err))
	}
	log.Info("minio客户端初始化完毕")
	return configs.MinioClient
}
