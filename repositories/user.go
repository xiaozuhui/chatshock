package repositories

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-28 14:25:14
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:46:43
 * @Description:
 */

import (
	"chatshock/configs"
	"chatshock/custom"
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

func (u UserRepo) FindUserByEmail(email string) (*entities.UserEntity, error) {
	var user models.UserModel
	var avatar models.FileModel
	err := configs.DBEngine.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	ent := user.ModelToEntity().(*entities.UserEntity)
	err = configs.DBEngine.First(&avatar, "uuid = ?", user.Avatar).Error
	if err != nil {
		return nil, err
	}
	ent.Avatar = avatar.ModelToEntity().(*entities.FileEntity)
	return ent, nil
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
	var avatar models.FileModel
	err := configs.DBEngine.First(&user, "uuid = ?", ID).Error
	if err != nil {
		return nil, err
	}
	ent := user.ModelToEntity().(*entities.UserEntity)
	err = configs.DBEngine.First(&avatar, "uuid = ?", user.Avatar).Error
	if err != nil {
		return nil, err
	}
	ent.Avatar = avatar.ModelToEntity().(*entities.FileEntity)
	return ent, nil
}

func (u UserRepo) FindUsers(IDs []uuid.UUID) ([]*entities.UserEntity, error) {
	var users []models.UserModel
	var avatars []models.FileModel

	iUsers := make([]custom.IModel, 0)
	res := make([]*entities.UserEntity, 0)

	err := configs.DBEngine.Where("uuid IN (?)", IDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	avatarUUIDs := make([]uuid.UUID, 0, 0)
	userDict := make(map[uuid.UUID]uuid.UUID, 0)          // UserUUID - AvatarUUID
	avatarDict := make(map[uuid.UUID]models.FileModel, 0) // AvatarUUID - Avatar
	for _, friend := range users {
		iUsers = append(iUsers, friend)
		avatarUUIDs = append(avatarUUIDs, friend.Avatar)
		userDict[friend.UUID] = friend.Avatar
	}
	err = configs.DBEngine.Where("uuid IN (?)", avatarUUIDs).Find(&avatars).Error
	if err != nil {
		return nil, err
	}
	for _, avatar := range avatars {
		avatarDict[avatar.UUID] = avatar
	}
	us := custom.DBs(iUsers)
	for _, u := range us {
		ue := u.(*entities.UserEntity)
		ue.Avatar = avatarDict[userDict[ue.UUID]].ModelToEntity().(*entities.FileEntity)
		res = append(res, ue)
	}
	return res, nil
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
	pass, err := utils.MakePassword(userEntity.UUID, userEntity.Password)
	if err != nil {
		return nil, err
	}
	userEntity.Password = pass
	// 转为model
	user, err := models.EntityToUserModel(&userEntity)
	if err != nil {
		return nil, err
	}
	baseModel, err := models.NewBaseModel(userEntity.UUID)
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
	if userEntity.Email != "" {
		userModel.Email = userEntity.Email
	}
	if userEntity.Password != "" {
		pass, err := utils.MakePassword(userEntity.UUID, userEntity.Password)
		if err != nil {
			return err
		}
		userModel.Password = pass
	}
	if userEntity.Avatar != nil {
		userModel.Avatar = userEntity.Avatar.UUID
	}
	err = configs.DBEngine.Model(&userModel).Save(&userModel).Error
	if err != nil {
		return err
	}
	return nil
}

var _ interfaces.IUser = UserRepo{}
