package models

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

//func (m ChatRoomUser) ModelToEntity() interface{} {
//	chatRoomUser := &entities.ChatRoomUser{}
//	chatRoomUser.UserID = m.UserID
//	chatRoomUser.EnterAt = m.EnterAt
//	chatRoomUser.LeaveAt = m.LeaveAt
//	chatRoomUser.IPAddress = m.IPAddress
//	return chatRoomUser
//}
//
//type ChatRoomUser struct {
//	UserID    uuid.UUID  `json:"uuid" gorm:"type:char(36);primary_key"`
//	EnterAt   time.Time  `json:"enter_at"` // 加入聊天室的时间
//	LeaveAt   *time.Time `json:"leave_at"` // 离开聊天室的时间
//	IPAddress string     `json:"ip_address"`
//}
