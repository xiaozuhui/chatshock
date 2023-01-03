package websockets

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-02 12:22:19
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-15 10:39:52
 * @Description:
 */

import (
	"chatshock/entities"
	"chatshock/services"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/pkg/errors"
)

type Message struct {
	FromID     uuid.UUID              `json:"from_id"`
	FromName   string                 `json:"from_name"`   // 哪个用户发送的消息
	MType      MType                  `json:"type"`        // 请求或是来源的消息类型
	MsgType    MessageType            `json:"msg_type"`    // 消息类型
	MsgContent string                 `json:"msg_content"` // 消息内容
	Files      []*entities.FileEntity `json:"files"`       // 文件数组
	MsgTime    time.Time              `json:"msg_time"`    // 消息创建的时间
	To         uuid.UUID              `json:"to"`          // 发送给(提及)哪些用户(如果是@，可以@多个)
	ToChatRoom uuid.UUID              `json:"to_chatroom"` // 发送到哪个聊天室
	At         map[uuid.UUID]string   `json:"at"`          // @了哪些用户
}

// MessageIn 从客户端传进来的数据，将会转为message
type MessageIn struct {
	Type       int                  `json:"type"`
	MsgType    int                  `json:"msg_type"`
	MsgContent string               `json:"msg_content"`
	Files      []uuid.UUID          `json:"files"`
	To         uuid.UUID            `json:"to"`
	ToChatroom uuid.UUID            `json:"to_chatroom"`
	At         map[uuid.UUID]string `json:"at"`
}

// MessageType 消息类型
type MessageType uint8

// MType 消息来源类型
type MType uint8

const (
	MtText         MessageType = iota // 1、文字消息
	MtPhoto                           // 2、图片消息
	MtDynamicPhoto                    // 3、动图消息（表情）
	MtVideo                           // 4、视频消息
	MtVoice                           // 5、语音消息
	MtFile                            // 6、文件消息
	MtCard                            // 7、复合消息（文字消息和图片组合在一起的消息类型）,卡片消息
	MtURL                             // 8、连接消息（安全连接，包括名片也属于连接）
	MtVisitingCard                    // 9、封装了一层的连接消息
)

const (
	MsgTypeNormal   MType = iota // 正常消息，私信
	MsgTypeChatRoom              // 聊天室消息，broadcast广播
	MsgTypeSystem                // 系统消息，正常消息，可以是通知
	MsgTypeError                 // 系统消息，错误消息
	MsgTypeUserList
)

// NewMessage 创建新数据，将
func NewMessage(user *User, params map[string]interface{}) (*Message, error) {
	log.Debug(fmt.Sprintf("参数为：%v", params))
	messageIn := MessageIn{}
	marshal, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &messageIn)
	if err != nil {
		return nil, err
	}
	mType := MType(messageIn.Type)
	if mType == MsgTypeNormal {
		// 如果是私信，那么To必须存在
		if messageIn.To == uuid.Nil {
			return nil, errors.WithStack(errors.New("参数错误：To数据不存在"))
		}
	} else if mType == MsgTypeChatRoom {
		// 如果是群聊，那么ToChatroom存在
		if messageIn.ToChatroom == uuid.Nil {
			return nil, errors.WithStack(errors.New("参数错误：ToChatRoom数据不存在"))
		}
	}
	// 获取到file的实体
	var files []*entities.FileEntity
	msgType := MessageType(messageIn.MsgType)
	if !(msgType == MtText || msgType == MtURL || msgType == MtVisitingCard) {
		fileService := services.FileFactory()
		files, err = fileService.GetFiles(messageIn.Files...)
		if err != nil {
			return nil, err
		}
	}
	message := &Message{
		FromID:     user.UserEntity.UUID,
		FromName:   user.UserEntity.NickName,
		MType:      mType,
		MsgType:    msgType,
		To:         messageIn.To,
		ToChatRoom: messageIn.ToChatroom,
		At:         messageIn.At,
		MsgContent: messageIn.MsgContent,
		Files:      files,
		MsgTime:    time.Now(),
	}
	return message, nil
}

func ErrorMessage(user *User, err error) *Message {
	message := &Message{
		FromID:     uuid.Nil,
		FromName:   "系统管理员",
		MType:      MsgTypeError,
		MsgType:    MtText,
		To:         user.UserEntity.UUID,
		MsgContent: errors.WithStack(err).Error(),
		Files:      nil,
		MsgTime:    time.Now(),
	}
	return message
}

func WelComeMessage(user *User, chatRoom *ChatRoom) *Message {
	return &Message{
		FromID:     uuid.Nil,
		FromName:   "系统管理员",
		MType:      MsgTypeSystem,
		MsgType:    MtText,
		To:         user.UserEntity.UUID,
		MsgContent: fmt.Sprintf("欢迎【%s】进入【%s】聊天室", user.UserEntity.NickName, chatRoom.CR.Name),
		Files:      nil,
		MsgTime:    time.Now(),
	}
}

func LeavingMessage(user *User, chatRoom *ChatRoom) *Message {
	return &Message{
		FromID:     uuid.Nil,
		FromName:   "系统管理员",
		MType:      MsgTypeSystem,
		MsgType:    MtText,
		To:         user.UserEntity.UUID,
		MsgContent: fmt.Sprintf("用户【%s】离开【%s】聊天室", user.UserEntity.NickName, chatRoom.CR.Name),
		Files:      nil,
		MsgTime:    time.Now(),
	}
}
