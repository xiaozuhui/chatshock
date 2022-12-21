package models

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-21 10:06:17
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-21 14:52:20
 * @Description:
 */

import (
	"chatshock/custom"
	"chatshock/entities"

	"github.com/gofrs/uuid"
)

type ChatroomModel struct {
	BaseModel
	Name        string                   `json:"name" gorm:"type:char(512)"`
	Description string                   `json:"description" gorm:"type:char(512)"`
	Users       map[uuid.UUID]*UserModel `json:"users" gorm:"-"`
	Master      *UserModel               `json:"master" gorm:"-"`
}

func (m ChatroomModel) ModelToEntity() interface{} {
	chatRoom := &entities.ChatRoom{}
	baseEntity := m.BaseModel.ModelToEntity()
	chatRoom.BaseEntity = *baseEntity
	chatRoom.Name = m.Name
	chatRoom.Description = m.Description
	chatRoom.Master = m.Master.ModelToEntity().(*entities.UserEntity)
	for userID, chatroom := range m.Users {
		chatRoom.Users[userID] = (*chatroom).ModelToEntity().(*entities.UserEntity)
	}
	return chatRoom
}

func EntityToChatroomModel(e *entities.ChatRoom) (*ChatroomModel, error) {
	m := &ChatroomModel{}
	m.BaseModel = *EntityToBaseModel(&e.BaseEntity)
	m.Name = e.Name
	m.Description = e.Description
	m.Master, _ = EntityToUserModel(e.Master)
	for _, u := range e.Users {
		m.Users[u.UUID], _ = EntityToUserModel(u)
	}
	return m, nil
}

var _ custom.IModel = ChatroomModel{}
