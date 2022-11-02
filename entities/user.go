package entities

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 10:11:01
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-02 09:29:35
 * @Description: 用户实体
 */

import "time"

type UserEntity struct {
	BaseEntity
	NickName     string        `json:"nickname"`
	Password     string        `json:"password"`
	PhoneNumber  string        `json:"phone_number"`
	Avatar       string        `json:"avatar"`
	Introduction string        `json:"introduction"`
	LastLogin    time.Time     `json:"last_login"`
	Gender       GenderTypeStr `json:"gender"`
}

type GenderType uint8
type GenderTypeStr string

const (
	Male GenderType = iota
	Female
	MaleStr   GenderTypeStr = "male"
	FemaleStr GenderTypeStr = "female"
)

func (t GenderType) ParseGenderType() GenderTypeStr {
	if t == Male {
		return MaleStr
	} else {
		return FemaleStr
	}
}

func (t GenderTypeStr) ParseGenderStr() GenderType {
	if t == MaleStr {
		return Male
	} else {
		return Female
	}
}
