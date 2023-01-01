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
	"fmt"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"reflect"
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
	To         map[uuid.UUID]string   `json:"to"`          // 发送给(提及)哪些用户(如果是@，可以@多个)
	ToChatRoom map[uuid.UUID]string   `json:"to_chatroom"` // 发送到哪个聊天室
}

// MessageType 消息类型
type MessageType uint8

// MType 消息来源类型
type MType uint8

// MessageTypeStr 消息类型，对应字符串
type MessageTypeStr string

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

	MtTextStr         MessageTypeStr = "text"
	MtPhotoStr        MessageTypeStr = "photo"
	MtDynamicPhotoStr MessageTypeStr = "dynamic_photo"
	MtVideoStr        MessageTypeStr = "video"
	MtVoiceStr        MessageTypeStr = "voice"
	MtFileStr         MessageTypeStr = "file"
	MtCardStr         MessageTypeStr = "text_and_photo"
	MtURLStr          MessageTypeStr = "url"
	MtVisitingCardStr MessageTypeStr = "visiting_card"
)
const (
	MsgTypeNormal   MType = iota // 正常消息，私信
	MsgTypeChatRoom              // 聊天室消息，broadcast广播
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
func NewMessage(user *User, params map[string]interface{}) (*Message, error) {
	/**
	params
	content: string，消息体本身
	to: []uuid.UUID
	type: int MsgTypeNormal|MsgTypeChatRoom
	msg_type: string
	*/
	// TODO 这部分需要优化，需要更加简明
	// TODO to_chatroom需要增加一下
	log.Info(fmt.Sprintf("参数为：%v", params))
	var (
		//value any
		err error
	)
	// 消息内容，可以为空
	//value, err = CheckMap(params, "msg_content", "string")
	//if err != nil {
	//	return nil, err
	//}
	var msgContent = params["msg_content"].(string)
	// 被指定的用户的uuid
	//value, err = CheckMap(params, "to", "[]uuid.UUID")
	//if err != nil {
	//	return nil, err
	//}
	var to = params["to"].([]interface{})
	// 消息传递类型
	//value, err = CheckMap(params, "type", "int")
	//if err != nil {
	//	return nil, err
	//}
	var mType = int(params["type"].(float64))
	// 消息类型
	//value, err = CheckMap(params, "msg_type", "string")
	//if err != nil {
	//	return nil, err
	//}
	var msgType = params["msg_type"].(string)
	// 文件的uuid
	//value, err = CheckMap(params, "files", "[]uuid.UUID", true)
	//if err != nil {
	//	return nil, err
	//}
	var fileUUIDs []uuid.UUID
	//if value == nil {
	//	fileUUIDs = []uuid.UUID{uuid.Nil}
	//} else {
	//	fileUUIDs = params["files"].([]uuid.UUID)
	//}
	// 消息时间
	//value, err = CheckMap(params, "msg_time", "time.Time")
	//if err != nil {
	//	return nil, err
	//}
	var msgTime, _ = time.Parse(params["msg_time"].(string), "2006-01-02 15:04:05")
	// 判断to中的人是否在broadcast中，因为所有能够发送数据的人，都应该在broadcast里
	toUsers := make(map[uuid.UUID]string, 0)
	for _, toID := range to {
		id, _ := uuid.FromString(toID.(string))
		if u, ok := BroadCaster.Users[id]; ok {
			toUsers[u.UserEntity.UUID] = u.UserEntity.NickName
		}
	}
	// 获取到file的实体
	fileService := services.FileFactory()
	files, err := fileService.GetFiles(fileUUIDs...)
	if err != nil {
		return nil, err
	}
	//
	message := &Message{
		FromID:     user.UserEntity.UUID,
		FromName:   user.UserEntity.NickName,
		MType:      MType(mType),
		MsgType:    MessageTypeStr(msgType).Parse(),
		To:         toUsers,
		MsgContent: msgContent,
		Files:      files,
		MsgTime:    msgTime,
	}
	return message, nil
}

func CheckMap(params map[string]any, checkKey, checkType string, noExcept ...bool) (any, error) {
	var (
		err   error
		ok    bool
		value any
	)

	if value, ok = params[checkKey]; ok {
		valueType := reflect.TypeOf(value).Kind().String()
		if valueType == checkType {
			err = nil
		} else {
			err = errors.WithStack(errors.New(fmt.Sprintf("params中[%s]数据[%v]的类型错误，不是%s，而是%s", checkKey, value, checkType, valueType)))
			value = nil
		}
	} else {
		value = nil
		err = errors.WithStack(errors.New(fmt.Sprintf("key=[%s]的数据不存在", checkKey)))
	}
	if err != nil {
		if len(noExcept) == 0 || !noExcept[0] {
			return nil, err
		}
		return nil, nil
	}
	return value, nil
}

func ErrorMessage(user *User, err error) *Message {
	message := &Message{
		FromID:     uuid.Nil,
		FromName:   "系统管理员",
		MType:      MsgTypeError,
		MsgType:    MtText,
		To:         map[uuid.UUID]string{user.UserEntity.UUID: user.UserEntity.NickName},
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
		To:         map[uuid.UUID]string{user.UserEntity.UUID: user.UserEntity.NickName},
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
		To:         map[uuid.UUID]string{user.UserEntity.UUID: user.UserEntity.NickName},
		MsgContent: fmt.Sprintf("用户【%s】离开【%s】聊天室", user.UserEntity.NickName, chatRoom.CR.Name),
		Files:      nil,
		MsgTime:    time.Now(),
	}
}
