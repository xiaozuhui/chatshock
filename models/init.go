package models

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 10:28:13
 * @Description:
 */

import (
	"chatshock/configs"
	log "github.com/sirupsen/logrus"
)

// InitModel 初始化模型
func InitModel() error {
	log.Info("初始化模型开始...")
	var migrates []interface{}
	migrates = append(migrates,
		FileModel{},
		UserModel{},
		FriendsModel{})
	err := configs.DBEngine.AutoMigrate(migrates...)
	if err != nil {
		return err
	}
	log.Info("初始化模型结束.")
	return nil
}
