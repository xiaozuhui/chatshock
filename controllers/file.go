package controllers

import (
	"chatshock/entities"
	"chatshock/middlewares"
	"chatshock/services"
	"chatshock/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

/**
文件操作的相关接口
*/

type FileController struct {
}

func (e *FileController) Router(engine *gin.Engine) {
	crH := engine.Group("/v1/file")
	crH.Use(middlewares.JWTAuth())
	crH.POST("", e.UploadFile)
	crH.POST("/multi", e.UploadFiles)
	crH.GET("")
	crH.GET("/multi")
}

func (e *FileController) UploadFile(c *gin.Context) {
	claims, ok := c.Get("claims")
	if ok != true {
		panic(errors.WithStack(errors.New("没有用户信息")))
	}
	userID := claims.(utils.UserClaims).UUID
	fileService := services.FileFactory()
	uploadFile, err := c.FormFile("file")
	if err != nil {
		panic(errors.WithStack(err))
	}
	file, err := utils.UploadFiles(userID.String(), uploadFile.Filename, uploadFile)
	if err != nil {
		return
	}
	contentType := uploadFile.Header.Get("Content-Type")
	fileType, err := entities.ContentType2FileType(contentType)
	if err != nil {
		panic(errors.WithStack(err))
	}
	saveFile, err := fileService.SaveFile(file, fileType, contentType)
	if err != nil {
		return
	}
	c.JSON(200, saveFile)
}

func (e *FileController) UploadFiles(c *gin.Context) {
	claims, ok := c.Get("claims")
	if ok != true {
		panic(errors.WithStack(errors.New("没有用户信息")))
	}
	userID := claims.(utils.UserClaims).UUID
	fileService := services.FileFactory()
	form, err := c.MultipartForm()
	if err != nil {
		panic(errors.WithStack(err))
	}
	fs, ok := form.File["files"]
	if !ok {
		panic(errors.WithStack(errors.New("参数错误")))
	}
	fileInfo := make([]map[string]interface{}, 0, 0)
	for _, f := range fs {
		uploadInfo, err := utils.UploadFiles(userID.String(), f.Filename, f)
		if err != nil {
			panic(errors.WithStack(err))
		}
		fileInfo = append(fileInfo, map[string]interface{}{
			"info":         uploadInfo,
			"file_type":    entities.ContentType2FileType(f.Header.Get("Content-Type")),
			"content_type": f.Header.Get("Content-Type"),
		})
	}
	files, err := fileService.SaveFiles(fileInfo)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, files)
}
