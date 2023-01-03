package websockets

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-19 09:01:30
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-02 01:55:04
 * @Description:
 */

import (
	"chatshock/entities"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ChatRoom struct {
	CR              *entities.ChatRoom `json:"chatroom"`         // 对应的聊天室
	EnteringChannel chan *User         `json:"entering_channel"` // 进入聊天室的消息
	LeavingChannel  chan *User         `json:"leaving_channel"`  // 离开聊天室消息
	MessageChannel  chan *Message      `json:"message_channel"`  // 消息

	canInChatRoom chan uuid.UUID
	isInChatRoom  chan bool

	ChatRoomLock sync.Mutex
}

// Listen
/**
 * @description: 在chatRoom被创建出来后，就应该打开监听程序
 * @return {*}
 * @author: xiaozuhui
 */
func (c *ChatRoom) Listen() {
	log.Infof("【%s】正在监听...", c.CR.Name)
	for {
		select {
		case user := <-c.EnteringChannel:
			if user == nil {
				continue
			}
			// 1、给新进用户发送欢迎消息；
			// 2、给聊天室的其他用户发送群消息，提示有新用户加入
			c.BroadCast(WelComeMessage(user, c))
		case user := <-c.LeavingChannel:
			if user == nil {
				continue
			}
			// 给群用户发送离开消息
			c.BroadCast(LeavingMessage(user, c))
		case msg := <-c.MessageChannel:
			if msg == nil {
				continue
			}
			// 给群用户发送消息
			c.BroadCast(msg)
		case userID := <-c.canInChatRoom:
			// 判断是否已经在聊天室中
			_, ok := c.CR.Users[userID]
			c.isInChatRoom <- ok
		}
	}
}

// CanEnterChatRoom 判断能否加入聊天室
func (c *ChatRoom) CanEnterChatRoom(userID uuid.UUID) bool {
	log.Infof("【%v】判断是否在聊天室中开始", userID)
	c.canInChatRoom <- userID
	log.Infof("【%v】判断是否在聊天室中开始", userID)
	return !(<-c.isInChatRoom)
}

func NewChatRoom(crEntity *entities.ChatRoom) ChatRoom {
	return ChatRoom{
		CR:              crEntity,
		EnteringChannel: make(chan *User, 0),
		LeavingChannel:  make(chan *User, 0),
		MessageChannel:  make(chan *Message, 0),
		canInChatRoom:   make(chan uuid.UUID, 0),
		isInChatRoom:    make(chan bool, 0),
	}
}

// BroadCast 向聊天室中的所有用户广播
func (c *ChatRoom) BroadCast(msg *Message) {
	userIDs := make([]uuid.UUID, 0, len(c.CR.Users))
	for k := range c.CR.Users {
		userIDs = append(userIDs, k)
	}
	for _, uID := range userIDs {
		if u, ok := BroadCaster.Users[uID]; ok {
			u.MessageChannel <- msg
		} else {
			// TODO 同样要保存到数据库
		}
	}
}

// UpdateChatRoomUser
/**
 * @description: 更新chatroom的用户结构，一般在删除用户或是添加用户的时候，要做这个
 * @param {*entities.ChatRoom} crEntity
 * @return {*}
 * @author: xiaozuhui
 */
func (c *ChatRoom) UpdateChatRoomUser(crEntity *entities.ChatRoom) {
	c.ChatRoomLock.Lock()
	c.CR = crEntity
	c.ChatRoomLock.Unlock()
}
