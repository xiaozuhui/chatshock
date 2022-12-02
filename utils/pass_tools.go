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
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

var AppSecret = []byte("6c9c48b1-1ef6dygg-91191a14-b9u14756")

type Salt []byte

// GenerateSalt
/**
 * @description: 自动生成salt
 * @param {string} phoneNumber
 * @return {*}
 * @author: xiaozuhui
 */
func GenerateSalt(phoneNumber string) Salt {
	var nums []int
	for _, n := range strings.Split(phoneNumber, "") {
		k, err := strconv.Atoi(n)
		if err != nil {
			k = 5
		}
		nums = append(nums, k)
	}
	res := make([]byte, 0, 8)
	n1 := nums[0] + nums[10]
	n2 := nums[1] + nums[9]
	n3 := nums[2] + nums[8]
	n4 := nums[3] + nums[7]
	n5 := nums[4] + nums[6]
	n6 := nums[5] + nums[5]
	n7 := n1 + n6
	n8 := n2 + n5
	k1 := AppSecret[n1]
	k2 := AppSecret[n2]
	k3 := AppSecret[n3]
	k4 := AppSecret[n4]
	k5 := AppSecret[n5]
	k6 := AppSecret[n6]
	k7 := AppSecret[n7]
	k8 := AppSecret[n8]
	res = append(res, k1, k2, k3, k4, k5, k6, k7, k8)
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
func MakePassword(phoneNumber, password string) (string, error) {
	salt := GenerateSalt(phoneNumber)
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
func CheckPassword(phoneNumber, password1, password2 string) (bool, error) {
	salt := GenerateSalt(phoneNumber)
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
