package websockets

import (
	"chatshock/entities"
	"context"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type User struct {
	UserEntity        *entities.UserEntity `json:"user"`                // 对应用户
	MessageChannel    chan *Message        `json:"message_channel"`     // 消息队列
	Conn              *websocket.Conn      `json:"conn"`                // 用户的websocket链接
	EnterChatRoomTime *time.Time           `json:"enter_chatroom_time"` // 加入聊天室的时间
	LeaveChatRoomTime *time.Time           `json:"leave_chatroom_time"` // 离开聊天室的时间
	IPAddress         string               `json:"ip_address"`          // 当时的用户IP地址
}

// TODO 系统管理员需要在init数据表时，写入数据库
//var System = &User{UserEntity: nil, MessageChannel: make(chan *Message)} // 系统默认用户

func NewUser(userID uuid.UUID, userEntity *entities.UserEntity, conn *websocket.Conn) *User {
	// 如果存在直接获取
	if user, ok := BroadCaster.Users[userID]; ok {
		return user
	}
	user := &User{
		UserEntity:     userEntity,
		MessageChannel: make(chan *Message),
		Conn:           conn,
	}
	UserLock.Lock()
	BroadCaster.UserLinks[userID] = conn
	BroadCaster.Users[userID] = user
	UserLock.Unlock()
	return user
}

// SendMessage 发送消息
func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		err := wsjson.Write(ctx, u.Conn, msg)
		if err != nil {
			log.Error(errors.WithStack(err))
		}
	}
}

// ReceiveMessage 从前端获取websocket数据
func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]any
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.Conn, &receiveMsg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			log.Error(errors.WithStack(err))
			u.MessageChannel <- ErrorMessage(u, err)
			return err
		}
		sendMsg, err := NewMessage(u, receiveMsg)
		if err != nil {
			log.Error(errors.WithStack(err))
			u.MessageChannel <- ErrorMessage(u, err)
			continue
		}
		// 判断类型，发送到私信还是发送到聊天室，通过broadcast发送
		switch sendMsg.MType {
		case MsgTypeNormal:
			// 私信
			if bU, ok := BroadCaster.Users[sendMsg.To]; ok {
				bU.MessageChannel <- sendMsg
			} else {
				// TODO 如果用户没有链接，那么需要存入数据库，到时候由客户端主动拉取数据
				continue
			}
		case MsgTypeChatRoom:
			// 聊天室消息
			if room, ok := BroadCaster.ChatRooms[sendMsg.ToChatRoom]; ok {
				room.MessageChannel <- sendMsg
			}
		}
	}
}
