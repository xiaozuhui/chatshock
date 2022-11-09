package repositories

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-28 14:25:14
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 13:35:42
 * @Description:
 */

import (
	"chatshock/configs"
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/models"
	"chatshock/utils"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type UserRepo struct {
}

// FindUser
/**
 * @description:  通过ID获取用户信息
 * @param {uuid.UUID} ID
 * @return {(*entities.UserEntity, error)}
 * @author: xiaozuhui
 */
func (u UserRepo) FindUser(ID uuid.UUID) (*entities.UserEntity, error) {
	var user models.UserModel
	err := configs.DBEngine.First(&user, "uuid = ?", ID).Error
	if err != nil {
		return nil, err
	}
	ent := user.ModelToEntity().(*entities.UserEntity)
	return ent, nil
}

func (u UserRepo) FindUsers(IDs []uuid.UUID) ([]*entities.UserEntity, error) {
	var users []models.UserModel
	iUsers := make([]models.IModel, 0)
	res := make([]*entities.UserEntity, 0)

	err := configs.DBEngine.Where("id IN (?)", IDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	for _, friend := range users {
		iUsers = append(iUsers, friend)
	}
	us := models.DBs(iUsers)
	for _, f := range us {
		res = append(res, f.(*entities.UserEntity))
	}
	return res, nil
}

// FindUserByPhoneNumber
/**
 * @description:  通过电话号码获取用户信息
 * @param {string} phoneNumber
 * @return {*}
 * @author: xiaozuhui
 */
func (u UserRepo) FindUserByPhoneNumber(phoneNumber string) (*entities.UserEntity, error) {
	var user models.UserModel
	err := configs.DBEngine.First(&user, "phone_number = ?", phoneNumber).Error
	if err != nil {
		return nil, err
	}
	ent := user.ModelToEntity().(*entities.UserEntity)
	return ent, nil
}

// DeleteUser
/**
 * @description: 删除用户（软删除）
 * @param {uuid.UUID} ID
 * @return {*}
 * @author: xiaozuhui
 */
func (u UserRepo) DeleteUser(ID uuid.UUID) error {
	var user models.UserModel
	err := configs.DBEngine.First(&user, "uuid = ?", ID).Error
	if err != nil {
		return err
	}
	err = configs.DBEngine.Delete(&user).Error
	return err
}

// CreateUser
/**
 * @description:  创建用户
 * @param {entities.UserEntity} userEntity
 * @param {string} password
 * @return {*}
 * @author: xiaozuhui
 */
func (u UserRepo) CreateUser(userEntity entities.UserEntity) (*entities.UserEntity, error) {
	// 先生成密码
	pass, err := utils.MakePassword(userEntity.PhoneNumber, userEntity.Password)
	if err != nil {
		return nil, err
	}
	userEntity.Password = pass
	// 转为model
	user := models.EntityToUserModel(&userEntity)
	baseModel, err := models.NewBaseModel()
	if err != nil {
		return nil, err
	}
	user.BaseModel = *baseModel
	user.LastLogin = time.Now()
	// 创建
	err = configs.DBEngine.Model(&models.UserModel{}).Create(user).Error
	if err != nil {
		return nil, err
	}
	e := user.ModelToEntity().(*entities.UserEntity)
	return e, nil
}

// UpdateLastLogin
/**
 * @description: 更新最后登录时间
 * @param {interface{}} ID
 * @return {*}
 * @author: xiaozuhui
 */
func (u UserRepo) UpdateLastLogin(ID uuid.UUID) error {
	err := configs.DBEngine.Model(&models.UserModel{}).
		Where("uuid = ?", ID).
		UpdateColumn("last_login", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateAccount
/**
 * @description: 更新用户信息
 * @param {entities.UserEntity} userEntity
 * @return {*}
 * @author: xiaozuhui
 */
func (u UserRepo) UpdateAccount(userEntity entities.UserEntity) error {
	userID := userEntity.UUID
	if userID == uuid.Nil {
		return errors.WithStack(errors.New("用户ID为空"))
	}
	var userModel models.UserModel
	err := configs.DBEngine.First(&userModel, "uuid = ?", userID).Error
	if err != nil {
		return err
	}
	userModel.UpdatedAt = time.Now()
	if userEntity.NickName != "" {
		userModel.NickName = userEntity.NickName
	}
	if userEntity.Gender != "" {
		userModel.Gender = userEntity.Gender.ParseGenderStr()
	}
	if userEntity.PhoneNumber != "" {
		userModel.PhoneNumber = userEntity.PhoneNumber
	}
	if userEntity.Password != "" {
		pass, err := utils.MakePassword(userEntity.PhoneNumber, userEntity.Password)
		if err != nil {
			return err
		}
		userModel.Password = pass
	}
	if userEntity.Avatar != "" {
		userModel.Avatar = userEntity.Avatar
	}
	err = configs.DBEngine.Model(&userModel).Save(&userModel).Error
	if err != nil {
		return err
	}
	return nil
}

var _ interfaces.IUser = UserRepo{}
