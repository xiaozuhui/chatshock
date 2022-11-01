package entities

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 10:12:00
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 10:19:34
 * @Description:
 */

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	UUID      uuid.UUID      `json:"uuid" `
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
