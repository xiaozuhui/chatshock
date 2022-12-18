package websockets

import (
	"github.com/gofrs/uuid"
	"nhooyr.io/websocket"
	"sync"
)

/**
BroadCast是用来进行传输消息的最底层

主要需要：
  - 维护所有用户的链接
  - 维护所有用户的信息
  - 维护消息的传递，每个用户都应该有自己的消息通道

使用饿汉式的单例模式，因为后续会频繁调用这个broadcast实例
*/

var UserLock sync.Mutex

type broadcast struct {
	UserLinks map[uuid.UUID]*websocket.Conn `json:"user_links"` // 所有登录的用户的链接
	Users     map[uuid.UUID]*User           `json:"users"`      // 登录的用户就将被加入这个broadcast
	ChatRooms map[uuid.UUID]*ChatRoom       `json:"chat_rooms"` // 被创建的聊天室

	MessageChannel chan *Message `json:"message_channel"` // 聊天消息
}

var BroadCaster = &broadcast{
	UserLinks:      make(map[uuid.UUID]*websocket.Conn, 0),
	Users:          make(map[uuid.UUID]*User, 0),
	ChatRooms:      make(map[uuid.UUID]*ChatRoom, 0),
	MessageChannel: make(chan *Message),
}
