package utils

import (
	"chatshock/configs"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/rakyll/magicmime"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"
	"time"
)

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
	contentType, err := GetFileContentType(src)
	if err != nil {
		contentType = ""
	}
	info, err := configs.MinioClient.PutObject(ctx,
		bucketName, objectName, src, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &info, nil
}

func GetFileUrl(bucketName, objectName string) (*url.URL, error) {
	ctx := context.Background()
	_url, err := configs.MinioClient.PresignedGetObject(ctx, bucketName, objectName, time.Second*24*60*60, url.Values{})
	if err != nil {
		return nil, err
	}
	return _url, nil
}

func GetFileContentType(file multipart.File) (string, error) {
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		log.Fatal(err)
	}
	defer magicmime.Close()
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	mimetype, err := magicmime.TypeByBuffer(fileByte)
	if err != nil {
		return "", err
	}
	return mimetype, nil
}
