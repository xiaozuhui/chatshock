package repositories

import (
	"chatshock/configs"
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/models"
	"github.com/gofrs/uuid"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-04 16:41:39
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-04 16:54:51
 * @Description:
 */

type FriendRepo struct {
}

func dbs(friends []models.FriendsModel) []*entities.FriendsEntity {
	var friendEntity = make([]*entities.FriendsEntity, 0, 0)
	for _, friend := range friends {
		friendEntity = append(friendEntity, friend.ModelToEntity())
	}
	return friendEntity
}

func (f FriendRepo) GetFriends(userID uuid.UUID) ([]*entities.FriendsEntity, error) {
	var friends []models.FriendsModel
	err := configs.DBEngine.First(&friends, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return dbs(friends), nil
}

// GetBindFriends 获取双向绑定的用户
func (f FriendRepo) GetBindFriends(userID uuid.UUID) ([]*entities.FriendsEntity, error) {
	var friends []models.FriendsModel
	err := configs.DBEngine.Raw(
		`SELECT f.* 
             FROM friends_model as f
			 LEFT JOIN friends_model as fm 
			 ON fm.user_id=f.other_id
				WHERE f.user_id=? AND f.deleted_at is null 
			    	AND fm.other_id=f.user_id AND f.deleted_at is null`,
		userID).Scan(&friends).Error
	if err != nil {
		return nil, err
	}
	return dbs(friends), nil
}

// GetUnBindFriends 获取这类用户，如果搜索某用户已经删掉了该用户的好友，但该用户还有此用户的好友，那么就搜索出来
func (f FriendRepo) GetUnBindFriends(userID uuid.UUID) ([]*entities.FriendsEntity, error) {
	var friends []models.FriendsModel
	err := configs.DBEngine.Raw(
		`SELECT f.* 
             FROM friends_model as f
			 LEFT JOIN friends_model as fm 
			 ON fm.user_id!=f.other_id
				WHERE f.user_id=? 
					AND f.deleted_at is null 
					AND fm.other_id=f.user_id 
					AND fm.deleted_at is not null`,
		userID).Scan(&friends).Error
	if err != nil {
		return nil, err
	}
	return dbs(friends), nil
}

// IsBindFriend 判断两个用户是不是互为好友
func (f FriendRepo) IsBindFriend(userID, otherID uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

// AddFriend 添加好友
func (f FriendRepo) AddFriend(userID, otherID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

// DeleteFriend 删除好友
func (f FriendRepo) DeleteFriend(userID, otherID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

var _ interfaces.IFriend = FriendRepo{}
