package models

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-03 15:31:47
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 12:44:08
 * @Description: 申请好友的时候，将会关联两者的uuid，双向绑定，两条数据；
 			     但是当其中一人删除好友的时候，只会保留没有删除的那人的数据；
				 不过后续要加回来的时候，如果对方还保留绑定，那么就不需要重新申请，可以直接加好友
 * @Tips: 必然存在两条数据，申请的时候就会同时创建两条数据，通过后即会生成两条数据
*/

import (
	"chatshock/entities"

	"github.com/gofrs/uuid"
)

// FriendsModel
/**
 * @description: 好友模型
 * @author: xiaozuhui
 */
type FriendsModel struct {
	BaseModel
	UserUUID  uuid.UUID `json:"user_id" gorm:"type:char(36)"`
	OtherUUID uuid.UUID `json:"other_id" gorm:"type:char(36)"`
}

func MakeFriendModel(userID, otherID uuid.UUID) (*FriendsModel, error) {
	BaseModel, err := NewBaseModel()
	if err != nil {
		return nil, err
	}
	fm := &FriendsModel{
		BaseModel: *BaseModel,
		UserUUID:  userID,
		OtherUUID: otherID,
	}
	return fm, nil
}

func (m FriendsModel) ModelToEntity() interface{} {
	fm := &entities.FriendsEntity{}
	fm.BaseEntity = *m.BaseModel.ModelToEntity()
	fm.UserUUID = m.UserUUID
	fm.OtherUUID = m.OtherUUID
	return fm
}

func EntityToFriendModel(e *entities.FriendsEntity) *FriendsModel {
	m := &FriendsModel{}
	m.BaseModel = *EntityToBaseModel(&e.BaseEntity)
	m.UserUUID = e.UserUUID
	m.OtherUUID = e.OtherUUID
	return m
}

var _ IModel = FriendsModel{}
