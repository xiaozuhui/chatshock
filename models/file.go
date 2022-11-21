package models

import (
	"chatshock/custom"
	"chatshock/entities"
	"github.com/gofrs/uuid"
	"time"
)

// FileModel 文件模型
/*
该模型不是业务模型，而是用来操作各种文件本身的模型

其中主要存储的是：
	1、bucket名称，一般是用户的手机号码
	2、文件名（存储在minio中的名称）
	3、可访问的url
	4、url失效时间
	5、文件类型：图片、视频、语音、动图
	6、mime类型
*/
type FileModel struct {
	UUID          uuid.UUID         `json:"uuid" gorm:"type:char(36);primary_key"`
	Bucket        string            `json:"bucket"`          // bucket名称
	FileName      string            `json:"filename"`        // 文件名称
	FileURL       string            `json:"file_url"`        // 文件可访问url
	URLExpireTime *time.Time        `json:"url_expire_time"` // url失效时间
	FileType      entities.FileType `json:"file_type"`       // 文件类型
	ContentType   string            `json:"mime_type"`       // MIME类型
}

func (m FileModel) ModelToEntity() interface{} {
	fm := &entities.FileEntity{}
	fm.UUID = m.UUID
	fm.FileName = m.FileName
	fm.Bucket = m.Bucket
	fm.FileURL = m.FileURL
	fm.URLExpireTime = m.URLExpireTime
	fm.FileType = m.FileType.Parse()
	fm.ContentType = m.ContentType
	return fm
}

func EntityToFileModel(e *entities.FileEntity) *FileModel {
	m := &FileModel{}
	m.UUID = e.UUID
	m.FileName = e.FileName
	m.Bucket = e.Bucket
	m.FileURL = e.FileURL
	m.URLExpireTime = e.URLExpireTime
	m.FileType = e.FileType.Parse()
	m.ContentType = e.ContentType
	return m
}

var _ custom.IModel = FileModel{}
