package models

import (
	"chatshock/configs"
	"fmt"
)

// InitModel 初始化模型
func InitModel() error {
	fmt.Println("初始化模型开始...")
	var migrates []interface{}
	migrates = append(migrates, User{})
	err := configs.DBEngine.AutoMigrate(migrates...)
	if err != nil {
		return err
	}
	fmt.Println("初始化模型结束.")
	return nil
}
