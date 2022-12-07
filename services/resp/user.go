package resp

import (
	"chatshock/entities"
	"chatshock/utils"
	"time"

	"github.com/gofrs/uuid"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-05 15:01:07
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:47:12
 * @Description:
 */

type User struct {
	UUID         uuid.UUID `json:"id"`
	NickName     string    `json:"nickname"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"`
	Introduction string    `json:"introduction"`
	Avatar       string    `json:"avatar"`
	LastLogin    time.Time `json:"last_login"`
}

// MakeUser
/**
 * @description: 构造用户返回值
 * @param {entities.UserEntity} userEntity
 * @return {*}
 * @author: xiaozuhui
 */
func MakeUser(userEntity entities.UserEntity) (*User, error) {
	user := User{
		UUID:         userEntity.UUID,
		NickName:     userEntity.NickName,
		PhoneNumber:  userEntity.PhoneNumber,
		Email:        userEntity.Email,
		Gender:       string(userEntity.Gender),
		LastLogin:    userEntity.LastLogin,
		Introduction: userEntity.Introduction,
	}
	// 如果没有或是过期了就从minio中获取
	if userEntity.Avatar == nil ||
		userEntity.Avatar.FileURL == "" ||
		time.Now().After(*userEntity.Avatar.URLExpireTime) {
		url, err := utils.GetFileUrl(userEntity.PhoneNumber, userEntity.PhoneNumber+"_avatar.png")
		if err != nil {
			return nil, err
		}
		user.Avatar = url.String()
	} else {
		user.Avatar = userEntity.Avatar.FileURL
	}
	return &user, nil
}
