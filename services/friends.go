package services

import (
	"chatshock/interfaces"
	"chatshock/repositories"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-08 21:37:53
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-08 21:47:57
 * @Description:
 */

type FriendService struct {
	friendRepo interfaces.IFriend
}

func FriendFactory() FriendService {
	fr := repositories.FriendRepo{}
	fs := FriendService{fr}
	return fs
}
