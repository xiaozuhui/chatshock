package applications

import (
	"chatshock/entities"
	"chatshock/services"
	"chatshock/utils"
	"errors"
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
	t := services.MakeToken(token, refresh, *expireTime)
	user, err := services.MakeUser(*ue)
	if err != nil {
		return nil, err
	}
	userResp := services.UserInfo{
		User:  user,
		Token: t,
	}
	return &userResp, nil
}
