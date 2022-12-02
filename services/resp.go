package services

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 15:53:19
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-06 09:39:30
 * @Description:
 */

import (
	"chatshock/services/resp"
)

// UserInfo
/**
 * @description: 包含User信息和Token信息的返回值
 * @author: xiaozuhui
 */
type UserInfo struct {
	User  *resp.User  `json:"user"`
	Token *resp.Token `json:"token"`
}
