package applications

import (
	"chatshock/services"

	"github.com/gofrs/uuid"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-10 09:52:56
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-10 10:47:57
 * @Description:
 */

type FriendApplication struct {
	UserService   services.UserService
	FriendService services.FriendService
}

func NewFriendApplication() FriendApplication {
	return FriendApplication{
		UserService:   services.UserFactory(),
		FriendService: services.FriendFactory(),
	}
}

// GetAllFriends
/**
 * @description: 获取所有的好友
 * @param {uuid.UUID} userID
 * @return {*}
 * @author: xiaozuhui
 */
func (a FriendApplication) GetAllFriends(userID uuid.UUID) ([]*services.User, error) {
	friends, err := a.FriendService.GetFriends(userID)
	if err != nil {
		return nil, err
	}
	otherUsers, err := a.UserService.GetUsers(friends)
	if err != nil {
		return nil, err
	}
	return otherUsers, nil
}

// GetFriends
/**
 * @description: 按照是否双向好友分类获取用户信息
 * @param {uuid.UUID} userID
 * @return {*}
 * @author: xiaozuhui
 */
func (a FriendApplication) GetFriends(userID uuid.UUID) (map[string][]*services.User, error) {
	bindUUIDs, unBindUUIDs, err := a.FriendService.GetBindFriends(userID)
	if err != nil {
		return nil, err
	}
	bindFriends, err := a.UserService.GetUsers(bindUUIDs)
	if err != nil {
		return nil, err
	}
	unBindFriends, err := a.UserService.GetUsers(unBindUUIDs)
	if err != nil {
		return nil, err
	}
	res := make(map[string][]*services.User, 0)
	res["bind"] = bindFriends
	res["unbind"] = unBindFriends
	return res, nil
}
