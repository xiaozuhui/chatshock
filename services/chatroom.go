package services

import (
	"chatshock/interfaces"
	"chatshock/repositories"
)

type ChatRoomService struct {
	ChatRoomRepo interfaces.IChatRoom
}

func ChatRoomFactory() ChatRoomService {
	crRepo := repositories.ChatRoomRepo{}
	crService := ChatRoomService{ChatRoomRepo: crRepo}
	return crService
}

func (s ChatRoomService) 