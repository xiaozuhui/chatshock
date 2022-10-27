/*
 * @Author: xiaozuhui
 * @Date: 2022-10-28 14:25:14
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 14:27:22
 * @Description:
 */
package models

import (
	"chatshock/entities"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	UUID      uuid.UUID      `json:"uuid" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func NewBaseModel() (*BaseModel, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	base := &BaseModel{
		UUID:      uuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return base, nil
}

func EntityToBaseModel(e *entities.BaseEntity) *BaseModel {
	m := &BaseModel{}
	m.CreatedAt = e.CreatedAt
	m.UUID = e.UUID
	m.DeletedAt = e.DeletedAt
	m.UpdatedAt = e.UpdatedAt
	return m
}

func (m *BaseModel) ModelToEntity() *entities.BaseEntity {
	e := &entities.BaseEntity{}
	e.UUID = m.UUID
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
	e.DeletedAt = m.DeletedAt
	return e
}
