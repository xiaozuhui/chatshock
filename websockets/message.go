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
}

// MessageType 消息类型
type MessageType uint8

// MType 消息来源类型
type MType uint8

// MessageTypeStr 消息类型，对应字符串
type MessageTypeStr string

const (
	MtText            MessageType    = iota // 1、文字消息
	MtPhoto                                 // 2、图片消息
	MtDynamicPhoto                          // 3、动图消息（表情）
	MtVideo                                 // 4、视频消息
	MtVoice                                 // 5、语音消息
	MtFile                                  // 6、文件消息
	MtCard                                  // 7、复合消息（文字消息和图片组合在一起的消息类型）,卡片消息
	MtURL                                   // 8、连接消息（安全连接，包括名片也属于连接）
	MtVisitingCard                          // 9、封装了一层的连接消息
	MtTextStr         MessageTypeStr = "text"
	MtPhotoStr        MessageTypeStr = "photo"
	MtDynamicPhotoStr MessageTypeStr = "dynamic_photo"
	MtVideoStr        MessageTypeStr = "video"
	MtVoiceStr        MessageTypeStr = "voice"
	MtFileStr         MessageTypeStr = "file"
	MtCardStr         MessageTypeStr = "text_and_photo"
	MtURLStr          MessageTypeStr = "url"
	MtVisitingCardStr MessageTypeStr = "visiting_card"
	MsgTypeNormal     MType          = iota // 正常消息，私信
	MsgTypeChatRoom                         // 聊天室消息，broadcast广播
	MsgTypeSystem                           // 系统消息，正常消息，可以是通知
	MsgTypeError                            // 系统消息，错误消息
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
	// TODO 创建要发送数据
	var content = params["msg_content"].(string)  // 消息内容，可以为空
	var to = params["to"].([]uuid.UUID)           // 被指定的用户的uuid
	var mType = params["type"].(int)              // 消息传递类型
	var msgType = params["msg_type"].(string)     // 消息类型
	var filesUUID = params["files"].([]uuid.UUID) // 文件模型的uuid
	var msgTime = params["msg_time"].(time.Time)  // 消息时间
	// 判断to中的人是否在broadcast中，因为所有能够发送数据的人，都应该在broadcast里
	toUsers := make(map[uuid.UUID]string, 0)
	for _, toID := range to {
		if u, ok := BroadCaster.Users[toID]; ok {
			toUsers[u.UserEntity.UUID] = u.UserEntity.NickName
		}
	}
	// 获取到file的实体
	fileService := services.FileFactory()
	files, err := fileService.GetFiles(filesUUID)
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
		MsgContent: content,
		Files:      files,
		MsgTime:    msgTime,
	}
	return message, nil
}
