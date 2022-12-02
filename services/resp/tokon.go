package resp

import "time"

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-05 15:01:16
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-05 15:31:22
 * @Description: Token相关的返回值
 */

type Token struct {
	Token      string    `json:"token"`       // token
	Refresh    string    `json:"refresh"`     // 刷新token
	ExpireTime time.Time `json:"expire_time"` // 过期时间
}

func MakeToken(token, refresh string, expireTime time.Time) *Token {
	tokenResp := &Token{
		Token:      token,
		Refresh:    refresh,
		ExpireTime: expireTime,
	}
	return tokenResp
}
