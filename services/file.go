package services

import (
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/repositories"
	"chatshock/utils"
	"github.com/gofrs/uuid"
	"github.com/minio/minio-go/v7"
	"time"
)

type FileService struct {
	FileRepo interfaces.IFile
}

func FileFactory() FileService {
	return FileService{
		FileRepo: repositories.FileRepo{},
	}
}

// GetFile 获取文件信息，主要是可访问URL
func (s FileService) GetFile(id uuid.UUID) (*entities.FileEntity, error) {
	// 获取文件信息
	file, err := s.FileRepo.GetFile(id)
	if err != nil {
		return nil, err
	}
	// 检查访问URL是否过期
	now := time.Now()
	if now.After(*file.URLExpireTime) {
		// 如果now大于过期时间，那么需要重新获取minio的url
		url_, err := utils.GetFileUrl(file.Bucket, file.FileName)
		if err != nil {
			return nil, err
		}
		// 保存到数据库
		file.FileURL = url_.String()
		et := file.URLExpireTime.Add(utils.Expires)
		file.URLExpireTime = &et
		file, err = s.FileRepo.UpdateFile(file.UUID, file)
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

// SaveFile 已经在方法外部，上传了文件
func (s FileService) SaveFile(info *minio.UploadInfo,
	fileType, mimeType string) (*entities.FileEntity, error) {
	file := &entities.FileEntity{}
	file.Bucket = info.Bucket
	file.FileName = info.Key
	file.FileType = entities.FileTypeStr(fileType)
	file.MIMEType = mimeType
	url_, err := utils.GetFileUrl(info.Bucket, info.Key)
	if err != nil {
		return nil, err
	}
	file.FileURL = url_.String()
	t := time.Now().Add(utils.Expires)
	file.URLExpireTime = &t
	id, err := s.FileRepo.SaveFile(file)
	if err != nil {
		return nil, err
	}
	getFile, err := s.FileRepo.GetFile(id)
	if err != nil {
		return nil, err
	}
	return getFile, nil
}

// SaveFiles 批量保存文件信息
func (s FileService) SaveFiles(infos []map[string]interface{}) ([]*entities.FileEntity, error) {
	fileEntities := make([]*entities.FileEntity, 0, 0)
	for _, info := range infos {
		uploadInfo := info["info"].(*minio.UploadInfo)
		fileType := info["file_type"].(string)
		mimeType := info["mime_type"].(string)
		file := &entities.FileEntity{}
		file.Bucket = uploadInfo.Bucket
		file.FileName = uploadInfo.Key
		file.FileType = entities.FileTypeStr(fileType)
		file.MIMEType = mimeType
		url_, err := utils.GetFileUrl(uploadInfo.Bucket, uploadInfo.Key)
		if err != nil {
			return nil, err
		}
		file.FileURL = url_.String()
		t := time.Now().Add(utils.Expires)
		file.URLExpireTime = &t
		fileEntities = append(fileEntities, file)
	}
	uuids, err := s.FileRepo.SaveFiles(fileEntities)
	if err != nil {
		return nil, err
	}
	files, err := s.FileRepo.GetFiles(uuids)
	if err != nil {
		return nil, err
	}
	return files, nil
}
