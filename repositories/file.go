package repositories

import (
	"chatshock/configs"
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/models"
	"fmt"
	"github.com/gofrs/uuid"
)

type FileRepo struct {
}

func (f FileRepo) SaveFile(fileEntity *entities.FileEntity) (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}
	fileEntity.UUID = id
	fm := models.EntityToFileModel(fileEntity)
	err = configs.DBEngine.Model(&models.FileModel{}).Create(fm).Error
	return id, err
}

func (f FileRepo) SaveFiles(fileEntities []*entities.FileEntity) ([]uuid.UUID, error) {
	fms := make([]*models.FileModel, 0, 0)
	ids := make([]uuid.UUID, 0, 0)
	for _, fileEntity := range fileEntities {
		id, err := uuid.NewV4()
		if err != nil {
			return ids, err
		}
		fileEntity.UUID = id
		fms = append(fms, models.EntityToFileModel(fileEntity))
		ids = append(ids, id)
	}
	err := configs.DBEngine.Model(&models.FileModel{}).CreateInBatches(fms, len(fms)).Error
	return ids, err
}

func (f FileRepo) GetFile(id uuid.UUID) (*entities.FileEntity, error) {
	var fm models.FileModel
	err := configs.DBEngine.First(&fm, "uuid=?", id).Error
	if err != nil {
		return nil, err
	}
	res := fm.ModelToEntity().(*entities.FileEntity)
	return res, nil
}

func (f FileRepo) GetFiles(ids []uuid.UUID) ([]*entities.FileEntity, error) {
	var fms []models.FileModel
	err := configs.DBEngine.Where("uuid in (?)", ids).Find(&fms).Error
	if err != nil {
		return nil, err
	}
	res := make([]*entities.FileEntity, 0, 0)
	for _, fm := range fms {
		res = append(res, fm.ModelToEntity().(*entities.FileEntity))
	}
	return res, nil
}

func (f FileRepo) DeleteFile(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (f FileRepo) DeleteFiles(ids []uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (f FileRepo) UpdateFile(id uuid.UUID, fileEntity *entities.FileEntity) (*entities.FileEntity, error) {
	fm := models.EntityToFileModel(fileEntity)
	err := configs.DBEngine.Model(fm).Updates(*fm).Error
	if err != nil {
		return nil, err
	}
	err = configs.DBEngine.First(fm, "uuid=?", id).Error
	if err != nil {
		return nil, err
	}
	res := fm.ModelToEntity().(*entities.FileEntity)
	return res, nil
}

// UpdateFiles 批量更新文件信息
func (f FileRepo) UpdateFiles(files map[uuid.UUID]*entities.FileEntity) ([]*entities.FileEntity, error) {
	var sql string
	ids := make([]uuid.UUID, 0, 0)
	var fms []models.FileModel

	for id, file := range files {
		sql += fmt.Sprintf(`UPDATE file_model 
							       SET file_url=%s, url_expire_time=%s 
 								   	   WHERE uuid=%v ;`, file.FileURL, file.URLExpireTime, id)
		ids = append(ids, id)
	}
	err := configs.DBEngine.Exec(sql).Error
	if err != nil {
		return nil, err
	}
	err = configs.DBEngine.Where("uuid IN (?)", ids).Find(&fms).Error
	if err != nil {
		return nil, err
	}
	res := make([]*entities.FileEntity, 0, 0)
	for _, fm := range fms {
		res = append(res, fm.ModelToEntity().(*entities.FileEntity))
	}
	return res, nil
}

var _ interfaces.IFile = FileRepo{}
