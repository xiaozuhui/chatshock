package services

import (
	"chatshock/interfaces"
	"chatshock/repositories"
	"chatshock/utils"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/gofrs/uuid"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-08 21:37:53
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 14:01:44
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

// GetFriends
/**
 * @description: 获取所有的好友
 * @param {uuid.UUID} userID
 * @return {[]uuid.UUID, error} 所有好友的uuid
 * @author: xiaozuhui
 */
func (s FriendService) GetFriends(userID uuid.UUID) ([]uuid.UUID, error) {
	friends, err := s.friendRepo.GetFriends(userID)
	if err != nil {
		return nil, err
	}
	otherIDs := make([]uuid.UUID, 0)
	for _, friend := range friends {
		otherIDs = append(otherIDs, friend.OtherUUID)
	}
	return otherIDs, nil
}

// GetBindFriends
/**
 * @description: 获取不同类型的好友
 * @param {uuid.UUID} userID
 * @return {[]uuid.UUID, []uuid.UUID, error} 1、是双向绑定的好友；2、是对方已经删除该用户的好友
 * @author: xiaozuhui
 */
func (s FriendService) GetBindFriends(userID uuid.UUID) ([]uuid.UUID, []uuid.UUID, error) {
	bindFriends, err := s.friendRepo.GetBindFriends(userID)
	if err != nil {
		return nil, nil, err
	}
	unBindFriends, err := s.friendRepo.GetUnBindFriends(userID)
	if err != nil {
		return nil, nil, err
	}
	bindUUID := make([]uuid.UUID, 0, len(bindFriends))
	unbindUUID := make([]uuid.UUID, 0, len(unBindFriends))
	for _, bindFriend := range bindFriends {
		bindUUID = append(bindUUID, bindFriend.OtherUUID)
	}
	for _, unBindFriend := range unBindFriends {
		unbindUUID = append(unbindUUID, unBindFriend.OtherUUID)
	}
	return bindUUID, unbindUUID, nil
}

// AddFriend
/**
* @description: 发送加好友申请
             	1、首先判断是否已经是好友
            	2、判断redis中是不是已经有相同的key
				3、在redis中塞入相关数据，保留7天
* @param {uuid.UUID} userID
* @param {uuid.UUID} otherID
* @return {error}
* @author: xiaozuhui
*/
func (s FriendService) AddFriend(userID, otherID uuid.UUID) error {
	// 1、首先判断是否已经是好友
	isBind, isFriend, isFriended, err := s.friendRepo.IsBindFriend(userID, otherID)
	if err != nil {
		return err
	}
	// 如果已经是好友，那么直接报错
	if isBind || isFriend {
		panic(errors.WithStack(errors.New("已经是好友，不需要重复申请")))
	}
	// 如果是对方好友，但不是自己好友，那么直接添加
	if isFriended {
		err := s.friendRepo.AddSideFriend(userID, otherID)
		if err != nil {
			return err
		}
		return nil
	}
	// 2、判断redis中是不是已经有相同的key
	get_, err := utils.RedisStrGet(fmt.Sprintf("%s-af-%s", userID, otherID))
	if err != nil {
		return err
	}
	if get_ != nil {
		panic(errors.WithStack(errors.New("已经发送申请，不需要重复申请")))
	}
	// 3、在redis中塞入相关数据
	expires := time.Second * 24 * 60 * 60 * 7 // 保留7天
	_, err = utils.RedisSet(fmt.Sprintf("%s-af-%s", userID, otherID), fmt.Sprintf("%s", otherID), &expires)
	if err != nil {
		return err
	}
	return nil
}

func (s FriendService) DeleteFriend(userID, otherID uuid.UUID) error {
	err := s.friendRepo.DeleteFriend(userID, otherID)
	return err
}

// ApplyFriend
/**
 * @description: 同意申请
 * @param {uuid.UUID} userID
 * @param {uuid.UUID} otherID
 * @return {error}
 * @author: xiaozuhui
 */
func (s FriendService) ApplyFriend(userID, otherID uuid.UUID) error {
	// 1、根据userID查找redis中的申请
	get_, err := utils.RedisStrGet(fmt.Sprintf("%s-af-%s", userID, otherID))
	if err != nil {
		return err
	}
	if get_ == nil {
		panic(errors.WithStack(errors.New("申请错误或是申请已过期")))
	}
	// 2、根据根据申请添加好友
	if !strings.EqualFold(*get_, otherID.String()) {
		panic(errors.WithStack(errors.New("申请错误")))
	}
	// 判断
	isBind, isFriend, isFriended, err := s.friendRepo.IsBindFriend(userID, otherID)
	if err != nil {
		return err
	}
	// 如果已经是好友，那么直接报错
	if isBind || isFriend {
		panic(errors.WithStack(errors.New("已经是好友，不需要重复申请")))
	}
	// 如果对方已经是好友，直接添加
	if isFriended {
		err := s.friendRepo.AddSideFriend(userID, otherID)
		if err != nil {
			return err
		}
		return nil
	}
	// 否则，同意申请，添加双方好友
	err = s.friendRepo.AddFriend(userID, otherID)
	return err
}
