package entities

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-15 08:33:28
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-15 13:09:26
 * @Description:
 */

import (
	"github.com/gofrs/uuid"
)

// ChatRoom 聊天室实体，需要固化到数据库
type ChatRoom struct {
	BaseEntity
	Name           string                    `json:"name"`        // 聊天室的名称
	Description    string                    `json:"description"` // 聊天室的介绍
	Users          map[uuid.UUID]*UserEntity `json:"users"`       // 聊天室中的用户
	Master         *UserEntity               `json:"master"`      // 聊天室的创始者
	ChatRoomAvatar *FileEntity               `json:"chatroom_avatar"`
}
