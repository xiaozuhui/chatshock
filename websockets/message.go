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
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Message struct {
	User        *entities.UserEntity   `json:"user"`         // 哪个用户发送的消息
	Type        MType                  `json:"type"`         // 请求或是来源的消息类型
	MessageType MessageType            `json:"message_type"` // 消息类型
	MsgContent  string                 `json:"msg_content"`  // 消息内容
	Files       []*entities.FileEntity `json:"files"`        // 文件数组
	SendTime    time.Time              `json:"send_time"`    // 发送时间
	MsgTime     time.Time              `json:"msg_time"`     // 消息创建的时间
}

type MessageType uint8
type MType uint8
type MessageTypeStr string

const (
	// 消息类型：
	MtText         MessageType = iota // 1、文字消息
	MtPhoto                           // 2、图片消息
	MtDynamicPhoto                    // 3、动图消息（表情）
	MtVideo                           // 4、视频消息
	MtVoice                           // 5、语音消息
	MtFile                            // 6、文件消息
	MtCard                            // 7、复合消息（文字消息和图片组合在一起的消息类型）,卡片消息
	MtURL                             // 8、连接消息（安全连接，包括名片也属于连接）
	MtVisitingCard                    // 9、封装了一层的连接消息
	// 消息类型，对应字符串
	MtTextStr         MessageTypeStr = "text"
	MtPhotoStr        MessageTypeStr = "photo"
	MtDynamicPhotoStr MessageTypeStr = "dynamic_photo"
	MtVideoStr        MessageTypeStr = "video"
	MtVoiceStr        MessageTypeStr = "voice"
	MtFileStr         MessageTypeStr = "file"
	MtCardStr         MessageTypeStr = "text_and_photo"
	MtURLStr          MessageTypeStr = "url"
	MtVisitingCardStr MessageTypeStr = "visiting_card"
	// 消息来源类型
	MsgTypeNormal   MType = iota // 正常消息，私信
	MsgTypeChatRoom              // 聊天室消息，broadcast广播
	MsgTypeGroup                 // 群组消息，群消息和聊天室的区别是，群需要有人拉进去或是申请，但是聊天室没有限制
	MsgTypeSystem                // 系统消息，正常消息，可以是通知
	MsgTypeError                 // 系统消息，错误消息
	MsgTypeUserList
)

// Parse 将类型MessageType转为MessageTypeStr类型
func (t MessageType) Parse() MessageTypeStr {
	var typ MessageTypeStr
	switch t {
	case MtText:
		typ = MtTextStr
	case MtPhoto:
		typ = MtPhotoStr
	case MtVideo:
		typ = MtVideoStr
	case MtVoice:
		typ = MtVoiceStr
	case MtFile:
		typ = MtFileStr
	case MtCard:
		typ = MtCardStr
	case MtDynamicPhoto:
		typ = MtDynamicPhotoStr
	case MtURL:
		typ = MtURLStr
	case MtVisitingCard:
		typ = MtVisitingCardStr
	default:
		panic(errors.WithStack(errors.New(fmt.Sprintf("错误的类型：%v", t))))
	}
	return typ
}

// Parse 将类型MessageTypeStr转为类型MessageType
func (t MessageTypeStr) Parse() MessageType {
	var typ MessageType
	switch t {
	case MtTextStr:
		typ = MtText
	case MtPhotoStr:
		typ = MtPhoto
	case MtVideoStr:
		typ = MtVideo
	case MtVoiceStr:
		typ = MtVoice
	case MtFileStr:
		typ = MtFile
	case MtCardStr:
		typ = MtCard
	case MtDynamicPhotoStr:
		typ = MtDynamicPhoto
	case MtURLStr:
		typ = MtURL
	case MtVisitingCardStr:
		typ = MtVisitingCard
	default:
		panic(errors.WithStack(errors.New(fmt.Sprintf("错误的类型：%v", t))))
	}
	return typ
}

// NewMessage 创建新数据，将
func NewMessage(user *User, params map[string]string) *Message {
	return nil
}
