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
	Name           string      `json:"name" gorm:"type:char(512);unique"` // 群名称
	Description    string      `json:"description" gorm:"type:char(512)"` // 备注/介绍
	Users          []uuid.UUID `json:"users"`                             // 群中用户
	Master         uuid.UUID   `json:"master" gorm:"type:char(36)"`       // 群主
	ChatRoomAvatar uuid.UUID   `json:"chatRoomAvatar"`                    // 聊天室头像
}

func (m ChatroomModel) ModelToEntity() interface{} {
	chatRoom := &entities.ChatRoom{}
	baseEntity := m.BaseModel.ModelToEntity()
	chatRoom.BaseEntity = *baseEntity
	chatRoom.Name = m.Name
	chatRoom.Description = m.Description
	return chatRoom
}

func EntityToChatroomModel(e *entities.ChatRoom) (*ChatroomModel, error) {
	m := &ChatroomModel{}
	m.BaseModel = *EntityToBaseModel(&e.BaseEntity)
	m.Name = e.Name
	m.Description = e.Description
	m.Master = e.Master.UUID
	for _, u := range e.Users {
		m.Users = append(m.Users, u.UUID)
	}
	m.ChatRoomAvatar = e.ChatRoomAvatar.UUID
	return m, nil
}

var _ custom.IModel = ChatroomModel{}
