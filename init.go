package main

import (
	"chatshock/configs"
	"chatshock/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// InitDatabase 初始化数据库
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

// InitConfig 初始化配置
func InitConfig() *configs.Config {
	configs.Conf = &configs.Config{}
	Pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
	}
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(Pwd + "/config/dev.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	configs.Conf.Parse(viper.GetViper())
	return configs.Conf
}

func InitRedis() *redis.Client {
	redisConfig := configs.Conf.RedisConfig
	configs.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.RedisHost, redisConfig.RedisPort),
		Password: "", // no password set
	})
	result := configs.RedisClient.Ping(context.Background())
	if result.Val() != "PONG" {
		// 连接有问题
		fmt.Println("redis链接不上")
		return nil
	}
	return configs.RedisClient
}

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
