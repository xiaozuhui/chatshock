package models

import (
	"chatshock/entities"
	"github.com/jackc/pgtype"
)

type MessageModel struct {
	BaseModel
	MessageType entities.MessageType `json:"message_type" gorm:"type:integer"` // 消息类型
	IsRead      bool                 `json:"is_read" gorm:"default:false"`     // 是否已读
	Content     string               `json:"content"`                          // 消息内容
	Files       pgtype.JSONArray     `json:"file_ids" gorm:"default:[]"`       // 文件的id
}
