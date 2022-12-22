package models

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-13 16:19:22
 * @Description:
 */

import (
	"chatshock/custom"
	"chatshock/entities"
	"errors"
	"time"
)

type UserModel struct {
	BaseModel
	NickName     string              `json:"nickname" gorm:"type:char(512);unique"` // nickname 即昵称
	Password     string              `json:"password" gorm:"type:char(512)"`        // 密码
	Email        string              `json:"email" gorm:"type:char(512);unique"`    // 邮箱地址，也可作为唯一标识
	Introduction string              `json:"introduction" gorm:"type:char(1024)"`   // 自我介绍
	Avatar       *FileModel          `json:"avatar" gorm:"-"`                       // 头像可能即存文件名称
	LastLogin    time.Time           `json:"last_login"`                            // 最后一次登录
	Gender       entities.GenderType `json:"gender" gorm:"type:integer"`            // 性别
}

func (m UserModel) ModelToEntity() interface{} {
	userEntity := &entities.UserEntity{}
	baseEntity := m.BaseModel.ModelToEntity()
	userEntity.BaseEntity = *baseEntity
	userEntity.NickName = m.NickName
	userEntity.Email = m.Email
	userEntity.LastLogin = m.LastLogin
	userEntity.Avatar = m.Avatar.ModelToEntity().(*entities.FileEntity)
	userEntity.Password = m.Password
	userEntity.Gender = m.Gender.ParseGenderType()
	userEntity.Introduction = m.Introduction
	return userEntity
}

func EntityToUserModel(e *entities.UserEntity) (*UserModel, error) {
	m := &UserModel{}
	if e.Email == "" {
		return nil, errors.New("电子邮箱不能为空")
	}
	m.BaseModel = *EntityToBaseModel(&e.BaseEntity)
	m.NickName = e.NickName
	m.Password = e.Password
	m.Email = e.Email
	m.LastLogin = e.LastLogin
	m.Avatar = EntityToFileModel(e.Avatar)
	m.Gender = e.Gender.ParseGenderStr()
	m.Introduction = e.Introduction
	return m, nil
}

var _ custom.IModel = UserModel{}
