package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	UUID      uuid.UUID      `json:"uuid" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
