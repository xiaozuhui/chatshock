package services

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 15:53:19
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 13:15:38
 * @Description:
 */

import (
	"chatshock/entities"
	"chatshock/utils"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Token struct {
	Token      string    `json:"token"`       // token
	Refresh    string    `json:"refresh"`     // 刷新token
	ExpireTime time.Time `json:"expire_time"` // 过期时间
}

type User struct {
	UUID         uuid.UUID `json:"id"`
	NickName     string    `json:"nickname"`
	PhoneNumber  string    `json:"phone_number"`
	Gender       string    `json:"gender"`
	Introduction string    `json:"introduction"`
	Avatar       string    `json:"avatar"`
	LastLogin    time.Time `json:"last_login"`
}

func MakeUser(userEntity entities.UserEntity) (*User, error) {
	user := User{
		UUID:         userEntity.UUID,
		NickName:     userEntity.NickName,
		PhoneNumber:  userEntity.PhoneNumber,
		Gender:       string(userEntity.Gender),
		LastLogin:    userEntity.LastLogin,
		Introduction: userEntity.Introduction,
	}
	// 查询redis中是否有对应的文件url
	avatarURL, err := utils.RedisStrGet(fmt.Sprintf("%s-avatar_url", userEntity.PhoneNumber))
	if err != nil {
		return nil, err
	}
	if avatarURL != nil {
		user.Avatar = *avatarURL
		return &user, nil
	}
	// 如果没有就从minio中获取
	url, err := utils.GetFileUrl(userEntity.PhoneNumber, userEntity.Avatar)
	if err != nil {
		return nil, err
	}
	user.Avatar = url.String()
	return &user, nil
}

func MakeToken(token, refresh string, expireTime time.Time) *Token {
	token_ := &Token{
		Token:      token,
		Refresh:    refresh,
		ExpireTime: expireTime,
	}
	return token_
}

type UserInfo struct {
	User  *User  `json:"user"`
	Token *Token `json:"token"`
}
