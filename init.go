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
	"log"
	"os"
	"path/filepath"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	smsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDatabase
/**
 * @description: 初始化数据库
 * @return {*}
 * @author: xiaozuhui
 */
func InitDatabase() {
	fmt.Println("初始化数据库开始...")
	var err error
	dbConfig := configs.Conf.DBConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBName, dbConfig.DBPort, dbConfig.SSLMode)
	configs.DBEngine, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}
	if configs.DBEngine == nil {
		fmt.Println("DB 为nil")
	}
	err = models.InitModel()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("初始化数据库结束.")
}

// InitConfig
/**
 * @description: 初始化配置，获取并解析configs中的配置文件
 * @param {string} configVersion
 * @return {*}
 * @author: xiaozuhui
 */
func InitConfig(configVersion string) *configs.Config {
	configs.Conf = &configs.Config{}
	configs.BaseDir, _ = os.Getwd()
	viper.SetConfigName(configVersion)
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(configs.BaseDir, "configs", configVersion+".yaml"))
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	configs.Conf.Parse(viper.GetViper())
	return configs.Conf
}

// InitRedis
/**
 * @description: 初始化redis
 * @return {*}
 * @author: xiaozuhui
 */
func InitRedis() *redis.Client {
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
	return configs.RedisClient
}

// InitMinioClient
/**
 * @description: 初始化minio，并获取其client
 * @return {*}
 * @author: xiaozuhui
 */
func InitMinioClient() *minio.Client {
	var err error
	configs.MinioClient, err = minio.New(configs.Conf.MinioConfig.EndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(configs.Conf.MinioConfig.AccessKeyID,
			configs.Conf.MinioConfig.SecretAccessKey, ""),
		Secure: configs.Conf.MinioConfig.UseSSL,
	})
	if err != nil {
		log.Fatal(errors.WithStack(err))
		return nil
	}
	return configs.MinioClient
}

// InitSmsClient
/**
 * @description:  初始化短信工具，并获取其SMSClient
 * @return {*}
 * @author: xiaozuhui
 */
func InitSmsClient() *smsapi.Client {
	var err error
	accessKeyId := tea.String(configs.Conf.PhoneConfig.AppKey)
	accessKeySecret := tea.String(configs.Conf.PhoneConfig.AppSecret)
	host := configs.Conf.PhoneConfig.Host
	_config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	_config.Endpoint = tea.String(host)
	configs.SMSClient, err = smsapi.NewClient(_config)
	if err != nil {
		log.Fatalf("%v", errors.WithStack(err))
	}
	return configs.SMSClient
}
