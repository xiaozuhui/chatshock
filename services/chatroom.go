package services

/*
 * @Author: xiaozuhui
 * @Date: 2023-01-01 22:30:51
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-02 01:43:52
 * @Description:
 */

import (
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/repositories"

	"github.com/gofrs/uuid"
)

type ChatRoomService struct {
	ChatRoomRepo interfaces.IChatRoom
}

func ChatRoomFactory() ChatRoomService {
	crRepo := repositories.ChatRoomRepo{}
	crService := ChatRoomService{ChatRoomRepo: crRepo}
	return crService
}

// GetChatRoom
/**
 * @description: 根据id获取聊天室
 * @param {...uuid.UUID} ids
 * @return {*}
 * @author: xiaozuhui
 */
func (s ChatRoomService) GetChatRoom(ids ...uuid.UUID) ([]*entities.ChatRoom, error) {
	chatRoomEntities := make([]*entities.ChatRoom, 0)
	if len(ids) == 1 {
		ent, err := s.ChatRoomRepo.FindChatRoom(ids[0])
		if err != nil {
			return nil, err
		}
		chatRoomEntities = append(chatRoomEntities, ent)
	}
	if len(ids) > 1 {
		ents, err := s.ChatRoomRepo.FindChatRooms(ids)
		if err != nil {
			return nil, err
		}
		chatRoomEntities = append(chatRoomEntities, ents...)
	}
	return chatRoomEntities, nil
}

// GetChatRoomsByUserIn
/**
 * @description: 搜索用户所在的聊天室
 * @param {uuid.UUID} userID
 * @return {*}
 * @author: xiaozuhui
 */
func (s ChatRoomService) GetChatRoomsByUserIn(userID uuid.UUID) ([]*entities.ChatRoom, error) {
	chatRoomEntities, err := s.ChatRoomRepo.FindChatRoomByUser(userID)
	return chatRoomEntities, err
}

func (s ChatRoomService) GetChatRoomByMaster(masterID uuid.UUID) ([]*entities.ChatRoom, error) {
	chatRoomEntities, err := s.ChatRoomRepo.FindChatRoomByMaster(masterID)
	return chatRoomEntities, err
}

// CreateChatRoom
/**
 * @description: 创建聊天室
 * @param {*entities.ChatRoom} chatRoomEntity
 * @return {*}
 * @author: xiaozuhui
 */
func (s ChatRoomService) CreateChatRoom(chatRoomEntity *entities.ChatRoom) (*entities.ChatRoom, error) {
	err := s.ChatRoomRepo.CreateChatRoom(chatRoomEntity)
	if err != nil {
		return nil, err
	}
	ent, err := s.ChatRoomRepo.FindChatRoom(chatRoomEntity.UUID)
	return ent, err
}

// IntoChatRoom
/**
 * @description: 加入某个聊天室
 * @param {*} userID
 * @param {uuid.UUID} chatroomID
 * @return {*}
 * @author: xiaozuhui
 */
func (s ChatRoomService) IntoChatRoom(userID, chatroomID uuid.UUID) (*entities.ChatRoom, error) {
	err := s.ChatRoomRepo.IntoChatRoom(userID, chatroomID)
	if err != nil {
		return nil, err
	}
	newEnt, err := s.ChatRoomRepo.FindChatRoom(chatroomID)
	return newEnt, err
}
