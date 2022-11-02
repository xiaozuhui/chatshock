package utils

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 12:22:19
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-02 09:39:06
 * @Description:
 */

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var CHARS = "ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789"

type ValidCodeType = int

const (
	RegisterOrLogin ValidCodeType = iota
)

func GetRandStr(n int) string {
	charArr := strings.Split(CHARS, "")
	charLen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))

	var res = ""
	for i := 0; i < n; i++ {
		res = res + charArr[ran.Intn(charLen)]
	}
	return res
}

type ValidCode struct {
	ValidCode  string        `json:"valid_code"`
	ExpireTime time.Duration `json:"expire_time"`
	CodeType   ValidCodeType `json:"code_type"`
}

// GenerateValidCode
/**
 * @description: 生成手机验证码
 * @param {ValidCodeType} t
 * @return {*}
 * @author: xiaozuhui
 */
func GenerateValidCode(t ValidCodeType) *ValidCode {
	v := ValidCode{}
	if t == RegisterOrLogin {
		v.CodeType = t
		v.ValidCode = v.registerCode()
		v.ExpireTime = time.Minute * 10
	}
	return &v
}

func (v *ValidCode) registerCode() string {
	return GetRandStr(6)
}

// CheckValidCode
/**
 * @description: 检查手机验证码
 * @param {string} phoneNumber 手机号
 * @param {string} vCode 验证码
 * @return {error} 如果正确，error为nil，否则存在错误
 * @author: xiaozuhui
 */
func CheckValidCode(phoneNumber string, vCode string) error {
	validCode, err := RedisStrGet(phoneNumber)
	if err != nil {
		log.Println(errors.WithStack(err))
		return err
	}
	if validCode == nil {
		return errors.WithStack(errors.New("验证码不存在或已经过期，请再次请求验证码"))
	}
	if strings.EqualFold(strings.ToUpper(vCode), strings.ToUpper(*validCode)) {
		return errors.WithStack(errors.New("验证码错误"))
	}
	return nil
}
