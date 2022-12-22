package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 10:52:31
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-21 16:16:55
 * @Description:
 */

import (
	"chatshock/configs"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisSet
/**
 * @description: 在redis中塞入数据已经过期时间
 * @param {*gin.Context} ctx
 * @param {string} k
 * @param {interface{}} v
 * @param {time.Duration} expTime
 * @return {*}
 * @author: xiaozuhui
 */
func RedisSet(k string, v interface{}, expTime *time.Duration) (string, error) {
	redisClient := configs.RedisClient
	ctx := context.Background()
	var result string
	var err error
	if expTime == nil {
		result, err = redisClient.Set(ctx, k, v, time.Minute*30).Result()
	} else {
		result, err = redisClient.Set(ctx, k, v, *expTime).Result()
	}

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return result, nil
}

func RedisDelete(k string) error {
	redisClient := configs.RedisClient
	ctx := context.Background()
	_, err := redisClient.Del(ctx, k).Result()
	return err
}

func RedisStrGet(k string) (*string, error) {
	redisClient := configs.RedisClient
	ctx := context.Background()
	res, err := redisClient.Get(ctx, k).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("键为 %s 的值不存在", k)
	}
	if err != nil {
		return nil, err
	}
	return &res, nil
}
