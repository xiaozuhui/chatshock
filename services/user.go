package services

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 15:20:26
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 15:06:16
 * @Description:
 */

import (
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/repositories"
	"chatshock/services/resp"
	"chatshock/utils"
	"errors"
	"gorm.io/gorm"

	"github.com/gofrs/uuid"
)

type UserService struct {
	UserRepo interfaces.IUser
}

func UserFactory() UserService {
	userRepo := repositories.UserRepo{}
	userService := UserService{userRepo}
	return userService
}

// CheckPassword
/**
 * @description: 检查密码是否正确
 * @param {string} phoneNumber
 * @param {string} password
 * @return {*}
 * @author: xiaozuhui
 */
func (s UserService) CheckPassword(contact interfaces.ISender, password string) (bool, error) {
	var user *entities.UserEntity
	var err error
	switch contact.Type() {
	case utils.Phone:
		user, err = s.UserRepo.FindUserByPhoneNumber(contact.String())
	case utils.Email:
		user, err = s.UserRepo.FindUserByEmail(contact.String())
	}
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("未注册")
	}
	isCheck, err := utils.CheckPassword(user.UUID, password, user.Password)
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
func (s UserService) Login(contact interfaces.ISender) (*UserInfo, error) {
	var userEntity *entities.UserEntity
	var err error
	switch contact.Type() {
	case utils.Phone:
		userEntity, err = s.UserRepo.FindUserByPhoneNumber(contact.String())
	case utils.Email:
		userEntity, err = s.UserRepo.FindUserByEmail(contact.String())
	}
	if err != nil {
		return nil, err
	}
	if userEntity == nil {
		return nil, errors.New("未注册")
	}
	token, refresh, expireTime, err := utils.GenerateToken(userEntity.UUID)
	if err != nil {
		return nil, err
	}
	// 更新user的login时间
	err = s.UserRepo.UpdateLastLogin(userEntity.UUID)
	if err != nil {
		return nil, err
	}
	user, err := resp.MakeUser(*userEntity)
	if err != nil {
		return nil, err
	}
	token_ := resp.MakeToken(token, refresh, *expireTime)
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
func (s UserService) GetUserByPhoneNumber(phoneNumber string) (*resp.User, error) {
	user, err := s.UserRepo.FindUserByPhoneNumber(phoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	userResp, err := resp.MakeUser(*user)
	if err != nil {
		return nil, err
	}
	return userResp, nil
}

func (s UserService) GetUserByEmailAddress(emailAddress string) (*resp.User, error) {
	user, err := s.UserRepo.FindUserByEmail(emailAddress)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	userResp, err := resp.MakeUser(*user)
	if err != nil {
		return nil, err
	}
	return userResp, nil
}

// GetUser
/**
 * @description: 通过ID获取用户
 * @param {uuid.UUID} userID
 * @return {*User} 用户信息
 * @author: xiaozuhui
 */
func (s UserService) GetUser(userID uuid.UUID) (*resp.User, error) {
	user, err := s.UserRepo.FindUser(userID)
	if err != nil {
		return nil, err
	}
	user_, err := resp.MakeUser(*user)
	if err != nil {
		return nil, err
	}
	return user_, nil
}

// GetUsers
/**
 * @description: 根据id批量获取用户
 * @param {[]uuid.UUID} userIDs
 * @return {*}
 * @author: xiaozuhui
 */
func (s UserService) GetUsers(userIDs []uuid.UUID) ([]*resp.User, error) {
	userEntities, err := s.UserRepo.FindUsers(userIDs)
	if err != nil {
		return nil, err
	}
	users := make([]*resp.User, 0, len(userEntities))
	for _, userEntity := range userEntities {
		user_, err := resp.MakeUser(*userEntity)
		if err != nil {
			return nil, err
		}
		users = append(users, user_)
	}
	return users, nil
}

// UpdateAccount
/**
 * @description: 更新账户信息
 * @param {*entities.UserEntity} userEntity
 * @return {*}
 * @author: xiaozuhui
 */
func (s UserService) UpdateAccount(userEntity *entities.UserEntity) error {
	err := s.UserRepo.UpdateAccount(*userEntity)
	return err
}
