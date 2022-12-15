package models

import (
	"chatshock/entities"
	"github.com/gofrs/uuid"
	"time"
)

type Chatroom struct {
	BaseModel
	Name        string `json:"name" gorm:"type:char(512)"`
	Description string `json:"description" gorm:"type:char(512)"`
	Users       map[uuid.UUID]*ChatRoomUser
	MasterID    uuid.UUID `json:"master_id"`
}

func (m Chatroom) ModelToEntity() interface{} {
	chatRoom := &entities.ChatRoom{}
	baseEntity := m.BaseModel.ModelToEntity()
	chatRoom.BaseEntity = *baseEntity
	chatRoom.Name = m.Name
	chatRoom.Description = m.Description
	chatRoom.Master = m.MasterID
	for userID, chatroom := range m.Users {
		chatRoom.Users[userID] = (*chatroom).ModelToEntity().(*entities.ChatRoomUser)
	}
	return chatRoom
}

func (m ChatRoomUser) ModelToEntity() interface{} {
	chatRoomUser := &entities.ChatRoomUser{}
	chatRoomUser.UserID = m.UserID
	chatRoomUser.EnterAt = m.EnterAt
	chatRoomUser.LeaveAt = m.LeaveAt
	chatRoomUser.IPAddress = m.IPAddress
	return chatRoomUser
}

type ChatRoomUser struct {
	UserID    uuid.UUID  `json:"uuid" gorm:"type:char(36);primary_key"`
	EnterAt   time.Time  `json:"enter_at"` // 加入聊天室的时间
	LeaveAt   *time.Time `json:"leave_at"` // 离开聊天室的时间
	IPAddress string     `json:"ip_address"`
}
