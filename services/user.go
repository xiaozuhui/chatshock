package services

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 15:20:26
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-04 11:49:40
 * @Description:
 */

import (
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/repositories"
	"chatshock/utils"
	"errors"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo interfaces.IUser
}

func UserFactory() UserService {
	userRepo := repositories.UserRepo{}
	userService := UserService{userRepo}
	return userService
}

// Register
/**
 * @description: 注册用户: 创建用户和账号信息
 * @param {entities.UserEntity} userEntity
 * @return {*}
 * @author: xiaozuhui
 */
func (s UserService) Register(userEntity entities.UserEntity) (*UserInfo, error) {
	if userEntity.PhoneNumber == "" {
		return nil, errors.New("手机号码不能为空")
	}
	_, err := s.userRepo.FindUserByPhoneNumber(userEntity.PhoneNumber)
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
	uploadInfo, err := utils.UploadImage(userEntity.PhoneNumber, userEntity.PhoneNumber+"_avatar.png", img)
	if err != nil {
		return nil, err
	}
	userEntity.Avatar = uploadInfo.Key
	ue, err := s.userRepo.CreateUser(userEntity)
	if err != nil {
		return nil, err
	}
	// 注册后默认登录
	token, refresh, expireTime, err := utils.GenerateToken(ue.PhoneNumber)
	if err != nil {
		return nil, err
	}
	t := MakeToken(token, refresh, *expireTime)
	user, err := MakeUser(*ue)
	if err != nil {
		return nil, err
	}
	userResp := UserInfo{
		User:  user,
		Token: t,
	}
	return &userResp, nil
}

// CheckPassword
/**
 * @description: 检查密码是否正确
 * @param {string} phoneNumber
 * @param {string} password
 * @return {*}
 * @author: xiaozuhui
 */
func (s UserService) CheckPassword(phoneNumber string, password string) (bool, error) {
	user, err := s.userRepo.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("该手机号没有注册")
	}
	isCheck, err := utils.CheckPassword(phoneNumber, password, user.Password)
	if err != nil {
		return false, err
	}
	return isCheck, nil
}

// Login
/**
 * @description: 登录
 * @param {string} phoneNumber
 * @return {UserInfo} 返回用户基本信息以及token
 * @author: xiaozuhui
 */
func (s UserService) Login(phoneNumber string) (*UserInfo, error) {
	userEntity, err := s.userRepo.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	if userEntity == nil {
		return nil, errors.New("该手机号没有注册")
	}
	token, refresh, expireTime, err := utils.GenerateToken(phoneNumber)
	if err != nil {
		return nil, err
	}
	// 更新user的login时间
	err = s.userRepo.UpdateLastLogin(userEntity.UUID)
	if err != nil {
		return nil, err
	}
	user, err := MakeUser(*userEntity)
	if err != nil {
		return nil, err
	}
	token_ := MakeToken(token, refresh, *expireTime)
	userInfo := &UserInfo{
		user, token_,
	}
	return userInfo, nil
}

// GetUserByPhoneNumber
/**
 * @description: 通过手机号获取用户
 * @param {string} phoneNumber
 * @return {*User, error}
 * @author: xiaozuhui
 */
func (s UserService) GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	user, err := s.userRepo.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("该手机号没有注册")
	}
	user_, err := MakeUser(*user)
	if err != nil {
		return nil, err
	}
	return user_, nil
}

// GetUser
/**
 * @description: 通过ID获取用户
 * @param {uuid.UUID} userID
 * @return {*User} 用户信息
 * @author: xiaozuhui
 */
func (s UserService) GetUser(userID uuid.UUID) (*User, error) {
	user, err := s.userRepo.FindUser(userID)
	if err != nil {
		return nil, err
	}
	user_, err := MakeUser(*user)
	if err != nil {
		return nil, err
	}
	return user_, nil
}

// UpdateAccount
/**
 * @description: 更新账户信息
 * @param {*entities.UserEntity} userEntity
 * @return {*}
 * @author: xiaozuhui
 */
func (s UserService) UpdateAccount(userEntity *entities.UserEntity) error {
	err := s.userRepo.UpdateAccount(*userEntity)
	return err
}
