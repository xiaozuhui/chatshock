package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 12:15:25
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 12:20:45
 * @Description:
 */

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/gofrs/uuid"

	"golang.org/x/crypto/scrypt"
)

//var AppSecret = []byte("6c9c48b1-1ef6dygg-91191a14-b9u14756")

type Salt []byte

// GenerateSalt
/**
 * @description: 自动生成salt
 * @param {string} phoneNumber
 * @return {*}
 * @author: xiaozuhui
 */
func GenerateSalt(UUID uuid.UUID) Salt {
	var nums []string
	nums = append(nums, strings.Split(strings.ReplaceAll(UUID.String(), "-", ""), "")...)
	n1 := nums[0] + nums[len(nums)-1]
	n2 := nums[1] + nums[len(nums)-2]
	n3 := nums[2] + nums[len(nums)-3]
	n4 := nums[3] + nums[len(nums)-4]
	n5 := nums[4] + nums[len(nums)-5]
	n6 := nums[5] + nums[len(nums)-6]
	n7 := n1 + n6
	n8 := n2 + n5
	res := []byte(n1 + n2 + n3 + n4 + n5 + n6 + n7 + n8)
	return res
}

// MakePassword
/**
 * @description: 生成Password
 * @param {*} phoneNumber
 * @param {string} password
 * @return {*}
 * @author: xiaozuhui
 */
func MakePassword(UUID uuid.UUID, password string) (string, error) {
	salt := GenerateSalt(UUID)
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dk), nil
}

// CheckPassword
/**
 * @description: 检查密码
 * @param {*} phoneNumber 用户手机号
 * @param {*} password1 用户填写密码
 * @param {string} password2 数据库保存密码
 * @return {*}
 * @author: xiaozuhui
 */
func CheckPassword(UUID uuid.UUID, password1, password2 string) (bool, error) {
	salt := GenerateSalt(UUID)
	dk, err := scrypt.Key([]byte(password1), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return false, err
	}
	pass := base64.StdEncoding.EncodeToString(dk)
	if pass != password2 {
		return false, errors.New("密码验证不正确")
	}
	return true, nil
}
