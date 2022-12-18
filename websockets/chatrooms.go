package websockets

import "chatshock/entities"

type ChatRoom struct {
	CR              *entities.ChatRoom `json:"chatroom"`         // 对应的聊天室
	EnteringChannel chan *User         `json:"entering_channel"` // 进入聊天室的消息
	LeavingChannel  chan *User         `json:"leaving_channel"`  // 离开聊天室消息
	MessageChannel  chan *Message      `json:"message_channel"`  // 消息
}
