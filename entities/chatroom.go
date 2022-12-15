package entities

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-15 08:33:28
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-15 13:09:26
 * @Description:
 */

import (
	"time"

	"github.com/gofrs/uuid"
	"nhooyr.io/websocket"
)

// ChatRoom 聊天室实体，需要固化到数据库
type ChatRoom struct {
	BaseEntity

	Name        string `json:"name"`        // 聊天室的名称
	Description string `json:"description"` // 聊天室的介绍

	Users  map[uuid.UUID]*ChatRoomUser `json:"users"`  // 聊天室中的用户
	Master uuid.UUID                   `json:"master"` // 聊天室的创始者

	EnteringChannel chan *ChatRoomUser `json:"entering_channel"` // 通知进入聊天室
	LeavingChannel  chan *ChatRoomUser `json:"leaving_channel"`  // 离开聊天室的通知
	MessageChannel  chan *Message      `json:"message_channel"`  // 聊天室消息

	//checkUserChannel      chan string
	//checkUserCanInChannel chan bool
}

// ChatRoomUser 聊天室和用户的映射关系，不需要固化到数据库
type ChatRoomUser struct {
	UserID         uuid.UUID       `json:"user_id"`
	EnterAt        time.Time       `json:"enter_at"` // 加入聊天室的时间
	LeaveAt        *time.Time      `json:"leave_at"` // 离开聊天室的时间
	IPAddress      string          `json:"ip_address"`
	MessageChannel chan *Message   `json:"message_channel"`
	Conn           *websocket.Conn `json:"conn"`
}

var Systemer = &ChatRoomUser{UserID: uuid.Nil}
