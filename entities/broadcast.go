package entities

import (
	"github.com/gofrs/uuid"
	"nhooyr.io/websocket"
	"time"
)

type Broadcast struct {
	Users  map[uuid.UUID]*BroadcastUser `json:"users"`
	Master uuid.UUID                    `json:"master"`

	EnteringChannel chan *BroadcastUser `json:"entering_channel"` // 通知进入聊天室
	LeavingChannel  chan *BroadcastUser `json:"leaving_channel"`  // 离开聊天室的通知
	MessageChannel  chan *Message       `json:"message_channel"`  // 聊天室消息

	checkUserChannel      chan string
	checkUserCanInChannel chan bool
}

type BroadcastUser struct {
	UID            uint64        `json:"uid"`
	UserID         uuid.UUID     `json:"user_id"`
	EnterAt        time.Time     `json:"enter_at"`
	IPAddress      string        `json:"ip_address"`
	MessageChannel chan *Message `json:"message_channel"`

	conn *websocket.Conn
}

var Systemer = &BroadcastUser{UID: 0, UserID: uuid.Nil}
