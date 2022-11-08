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

// dbs
/**
 * @description: 将多个model批量转为entity
 * @param {[]models.FriendsModel} friends
 * @return {[]*entities.FriendsEntity}
 * @author: xiaozuhui
 */
func dbs(friends []models.FriendsModel) []*entities.FriendsEntity {
	var friendEntity = make([]*entities.FriendsEntity, 0)
	for _, friend := range friends {
		friendEntity = append(friendEntity, friend.ModelToEntity())
	}
	return friendEntity
}

// GetFriends
/**
 * @description: 获取所有好友
 * @param {uuid.UUID} userID
 * @return {([]*entities.FriendsEntity, error)}
 * @author: xiaozuhui
 */
func (f FriendRepo) GetFriends(userID uuid.UUID) ([]*entities.FriendsEntity, error) {
	var friends []models.FriendsModel
	err := configs.DBEngine.Where("user_id = ?", userID).Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return dbs(friends), nil
}

// GetBindFriends
/**
 * @description: 获取双向绑定的用户
 * @param {uuid.UUID} userID
 * @return {([]*entities.FriendsEntity, error)}
 * @author: xiaozuhui
 */
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

// GetUnBindFriends
/**
 * @description: 获取这类用户，如果搜索某用户已经删掉了该用户的好友，但该用户还有此用户的好友，那么就搜索出来
 * @param {uuid.UUID} userID
 * @return {([]*entities.FriendsEntity, error)}
 * @author: xiaozuhui
 */
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

// IsBindFriend
/**
 * @description: 判断两个用户是不是互为好友
 * @param {uuid.UUID} userID
 * @param {uuid.UUID} otherID
 * @return {(bool, error)}
 * @author: xiaozuhui
 */
func (f FriendRepo) IsBindFriend(userID, otherID uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

// AddFriend
/**
 * @description: 添加好友
 * @param {uuid.UUID} userID
 * @param {uuid.UUID} otherID
 * @return {error}
 * @author: xiaozuhui
 */
func (f FriendRepo) AddFriend(userID, otherID uuid.UUID) error {
	fm1, err := models.MakeFriendModel(userID, otherID)
	if err != nil {
		return err
	}
	fm2, err := models.MakeFriendModel(otherID, userID)
	if err != nil {
		return err
	}
	err = configs.DBEngine.Model(&models.FriendsModel{}).Create(fm1).Error
	if err != nil {
		return err
	}
	err = configs.DBEngine.Model(&models.FriendsModel{}).Create(fm2).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteFriend
/**
 * @description: 删除好友
 * @param {uuid.UUID} userID
 * @param {uuid.UUID} otherID
 * @return {error}
 * @author: xiaozuhui
 */
func (f FriendRepo) DeleteFriend(userID, otherID uuid.UUID) error {
	friends := []models.FriendsModel{}
	err := configs.DBEngine.Where("user_id = ?", userID).Find(&friends).Error
	if err != nil {
		return err
	}
	for _, f := range friends {
		if f.OtherUUID == otherID {
			err = configs.DBEngine.Delete(&f).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var _ interfaces.IFriend = FriendRepo{}
