/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:17:18
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 10:28:13
 * @Description:
 */
package models

import (
	"chatshock/configs"
	"fmt"
)

// InitModel 初始化模型
func InitModel() error {
	fmt.Println("初始化模型开始...")
	var migrates []interface{}
	migrates = append(migrates, UserModel{})
	err := configs.DBEngine.AutoMigrate(migrates...)
	if err != nil {
		return err
	}
	fmt.Println("初始化模型结束.")
	return nil
}
