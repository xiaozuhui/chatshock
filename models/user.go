package models

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-04 13:35:13
 * @Description:
 */

import (
	"chatshock/entities"
	"time"
)

type UserModel struct {
	BaseModel
	NickName     string              `json:"nickname" gorm:"type:char(512)"`                          // nickname 即昵称
	Password     string              `json:"password" gorm:"type:char(512)"`                          // 密码
	PhoneNumber  string              `json:"phone_number" gorm:"type:char(11);unique_index;not null"` // 手机号码将作为唯一标识
	Avatar       string              `json:"avatar" gorm:"type:char(512)"`                            // 头像可能即存文件名称
	Introduction string              `json:"introduction" gorm:"type:char(1024)"`                     // 自我介绍
	LastLogin    time.Time           `json:"last_login"`                                              // 最后一次登录
	Gender       entities.GenderType `json:"gender" gorm:"type:char(50)"`                             // 性别
}

func (m *UserModel) ModelToEntity() *entities.UserEntity {
	userEntity := &entities.UserEntity{}
	baseEntity := m.BaseModel.ModelToEntity()
	userEntity.BaseEntity = *baseEntity
	userEntity.NickName = m.NickName
	userEntity.PhoneNumber = m.PhoneNumber
	userEntity.LastLogin = m.LastLogin
	userEntity.Avatar = m.Avatar
	userEntity.Password = m.Password
	userEntity.Gender = m.Gender.ParseGenderType()
	userEntity.Introduction = m.Introduction
	return userEntity
}

func EntityToUserModel(e *entities.UserEntity) *UserModel {
	m := &UserModel{}
	m.BaseModel = *EntityToBaseModel(&e.BaseEntity)
	m.NickName = e.NickName
	m.Password = e.Password
	m.PhoneNumber = e.PhoneNumber
	m.LastLogin = e.LastLogin
	m.Avatar = e.Avatar
	m.Gender = e.Gender.ParseGenderStr()
	m.Introduction = e.Introduction
	return m
}
