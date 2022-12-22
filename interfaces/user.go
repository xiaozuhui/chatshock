package interfaces

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 11:11:11
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 13:38:58
 * @Description:
 */

import (
	"chatshock/entities"

	"github.com/gofrs/uuid"
)

type IUser interface {
	// FindUser 根据id获取员工
	FindUser(ID uuid.UUID) (*entities.UserEntity, error)
	// FindUsers 根据ids获取员工
	FindUsers(IDs []uuid.UUID) ([]*entities.UserEntity, error)
	// FindUserByEmail 根据Email获取用户
	FindUserByEmail(email string) (*entities.UserEntity, error)
	// DeleteUser 删除账号
	DeleteUser(ID uuid.UUID) error
	// CreateUser 创建账号
	CreateUser(userEntity entities.UserEntity) (*entities.UserEntity, error)
	// UpdateLastLogin 更新最新的登录时间
	UpdateLastLogin(ID uuid.UUID) error
	// UpdateAccount 更新账户信息
	UpdateAccount(userEntity entities.UserEntity) error
}
