package applications

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-02 12:22:19
 * @LastEditors: xiaozuhui xiaozuhui@outlook.com
 * @LastEditTime: 2022-12-06 22:29:32
 * @Description:
 */

import (
	"chatshock/entities"
	"chatshock/services"
	"chatshock/services/resp"
	"chatshock/utils"
	"errors"
	"mime/multipart"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserApplication struct {
	UserService services.UserService
	FileService services.FileService
}

func NewUserApplication() UserApplication {
	return UserApplication{
		UserService: services.UserFactory(),
		FileService: services.FileFactory(),
	}
}

// Register
/**
 * @description: 注册用户: 创建用户和账号信息
 * @param {entities.UserEntity} userEntity
 * @return {*}
 * @author: xiaozuhui
 */
func (a UserApplication) Register(userEntity entities.UserEntity) (*services.UserInfo, error) {
	if userEntity.PhoneNumber == "" {
		return nil, errors.New("手机号码不能为空")
	}
	_, err := a.UserService.UserRepo.FindUserByPhoneNumber(userEntity.PhoneNumber)
	if err != gorm.ErrRecordNotFound {
		return nil, errors.New("该手机号码已经被注册")
	}
	// 创建默认头像
	img, err := utils.GenerateAvatar(userEntity.PhoneNumber)
	if err != nil {
		return nil, err
	}
	err = utils.MakeBucket(userEntity.PhoneNumber)
	if err != nil {
		return nil, err
	}
	// 保存文件
	uploadInfo, err := utils.UploadImage(userEntity.PhoneNumber, userEntity.PhoneNumber+"_avatar.png", img)
	if err != nil {
		return nil, err
	}
	// 保存文件信息
	fileEntity, err := a.FileService.SaveFile(uploadInfo, string(entities.PhotoStr), "application/png")
	if err != nil {
		return nil, err
	}
	userEntity.Avatar = fileEntity
	ue, err := a.UserService.UserRepo.CreateUser(userEntity)
	if err != nil {
		return nil, err
	}
	// 注册后默认登录
	token, refresh, expireTime, err := utils.GenerateToken(ue.PhoneNumber)
	if err != nil {
		return nil, err
	}
	t := resp.MakeToken(token, refresh, *expireTime)
	user, err := resp.MakeUser(*ue)
	if err != nil {
		return nil, err
	}
	userResp := services.UserInfo{
		User:  user,
		Token: t,
	}
	return &userResp, nil
}

// UpdateAvatar
/**
 * @description: 更新头像
 * @param {uuid.UUID} userID
 * @param {*multipart.FileHeader} avatar
 * @return (*resp.User, error)
 */
func (a UserApplication) UpdateAvatar(userID uuid.UUID, avatar *multipart.FileHeader) (*resp.User, error) {
	userResp, err := a.UserService.GetUser(userID)
	if err != nil {
		return nil, err
	}
	imgInfo, err := utils.UploadFiles(userResp.PhoneNumber, avatar.Filename, avatar)
	if err != nil {
		return nil, err
	}
	// 保存信息
	fileEntity, err := a.FileService.SaveFile(imgInfo, "photo", avatar.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	userEntity := entities.UserEntity{
		BaseEntity: entities.BaseEntity{
			UUID: userResp.UUID,
		},
		PhoneNumber: userResp.PhoneNumber,
		Avatar:      fileEntity,
	}
	err = a.UserService.UpdateAccount(&userEntity)
	if err != nil {
		return nil, err
	}
	userResp, err = a.UserService.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return userResp, nil
}
