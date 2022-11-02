package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-01 09:55:24
 * @Description:
 */

import (
	"chatshock/configs"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
)

// MakeBucket
/**
 * @description: 创建Minio桶
 * @param {string} bucketName 桶的名称，注意bucket一般为该用户的手机号码
 * @return {error} 如果有错误则返回错误，如果没有错误，则返回nil
 * @author: xiaozuhui
 */
func MakeBucket(bucketName string) error {
	ctx := context.Background()
	location := "us-east-1"
	err := configs.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := configs.MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists != nil || !exists {
			return err
		}
	}
	return nil
}

// UploadFiles
/**
 * @description: 上传文件
 * @param {string} bucketName 桶的名称
 * @param {string} objectName 文件的名称
 * @param {*multipart.FileHeader} file 文件流
 * @return {*minio.UploadInfo, error} 文件基础信息和错误
 * @author: xiaozuhui
 */
func UploadFiles(bucketName, objectName string, file *multipart.FileHeader) (*minio.UploadInfo, error) {
	ctx := context.Background()

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}(src)
	// Content-Type如果没有这个参数，则会返回空字符串
	contentType := file.Header.Get("Content-Type")
	// ContentType为空字符串的话所有的文件在访问到的时候就会下载
	info, err := configs.MinioClient.PutObject(ctx,
		bucketName, objectName, src, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// UploadImage
/**
 * @description:  上传PNG图片
 * @param {string} bucketName
 * @param {string} objectName
 * @param {*os.File} img
 * @return {(*minio.UploadInfo, error)}
 * @author: xiaozuhui
 */
func UploadImage(bucketName, objectName string, img *os.File) (*minio.UploadInfo, error) {
	ctx := context.Background()
	info, err := configs.MinioClient.PutObject(ctx,
		bucketName, objectName, img, -1, minio.PutObjectOptions{ContentType: "application/png"})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// GetFileUrl
/**
 * @description: 获取下载链接，这个下载链接是有过期时间的，所以每次请求都会重新生成
 * @tips: 在前端可以维护一个链接和资源的二维表，过期时间之前都可以直接用该链接，减少对资源的请求，或者在后端维护
 * @param {string} bucketName 桶名称，这里是用户的手机号码
 * @param {string} objectName 文件名称
 * @return {*url.URL, error} 文件可下载URL或错误
 * @author: xiaozuhui
 */
func GetFileUrl(bucketName, objectName string) (*url.URL, error) {
	ctx := context.Background()
	expires := time.Second * 24 * 60 * 60
	_url, err := configs.MinioClient.PresignedGetObject(ctx, bucketName, objectName, expires, url.Values{})
	if err != nil {
		return nil, err
	}
	// 取出url后，存入redis
	_, err = RedisSet(fmt.Sprintf("%s-avatar_url", bucketName), _url.String(), &expires)
	if err != nil {
		return nil, err
	}
	return _url, nil
}
