package interfaces

import (
	"chatshock/entities"
	"github.com/gofrs/uuid"
)

type IChatRoom interface {
	// FindChatRoom 获取聊天室
	FindChatRoom(ID uuid.UUID) (*entities.ChatRoom, error)
	// FindChatRoomByMaster 根据群主搜索聊天室
	FindChatRoomByMaster(masterID uuid.UUID) (*entities.ChatRoom, error)
	// CreateChatRoom 创建聊天室
	CreateChatRoom(masterID uuid.UUID) (*entities.ChatRoom, error)
	// DeleteChatRoom 删除聊天室
	DeleteChatRoom(ID uuid.UUID) error
	// DeleteChatRoomsByMaster 删除群主的所有聊天室
	DeleteChatRoomsByMaster(masterID uuid.UUID) error
	// FindChatRoomByUser 根据进入聊天室的用户获取所有存在这个用户的聊天室
	FindChatRoomByUser(userID uuid.UUID) ([]*entities.ChatRoom, error)
}
