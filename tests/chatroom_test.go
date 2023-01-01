package tests

import (
	"chatshock/configs"
	"chatshock/entities"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateChatRoom(t *testing.T) {
	configs.Conf = &configs.Config{}
	configs.BaseDir = "../"
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(configs.BaseDir, "configs", "dev.yaml"))
	err := viper.ReadInConfig()
	if err != nil {
		t.Error(err.Error())
	}
	configs.Conf.Parse(viper.GetViper())

	userMap := make(map[uuid.UUID]*entities.UserEntity, 0)
	id1, err := uuid.NewV4()
	if err != nil {
		return
	}
	id2, err := uuid.NewV4()
	if err != nil {
		return
	}
	id3, err := uuid.NewV4()
	if err != nil {
		return
	}
	userMap[id1] = &entities.UserEntity{
		BaseEntity:   *entities.NewBaseEntity(id1),
		NickName:     "",
		Password:     "",
		Email:        "",
		Avatar:       nil,
		Introduction: "",
		LastLogin:    time.Time{},
		Gender:       "",
	}
	userMap[id2] = &entities.UserEntity{
		BaseEntity:   *entities.NewBaseEntity(id2),
		NickName:     "",
		Password:     "",
		Email:        "",
		Avatar:       nil,
		Introduction: "",
		LastLogin:    time.Time{},
		Gender:       "",
	}
	userMap[id3] = &entities.UserEntity{
		BaseEntity:   *entities.NewBaseEntity(id3),
		NickName:     "",
		Password:     "",
		Email:        "",
		Avatar:       nil,
		Introduction: "",
		LastLogin:    time.Time{},
		Gender:       "",
	}
	chatRoomEntity := &entities.ChatRoom{
		BaseEntity:  *entities.NewBaseEntity(uuid.Nil),
		Name:        "test",
		Description: "用于测试",
		Users:       userMap,
		Master:      userMap[id1],
	}
	//crRepo := repositories.ChatRoomRepo{}
	//err = crRepo.CreateChatRoom(chatRoomEntity)
	//if err != nil {
	//	t.Error(err)
	//}
	t.Log(chatRoomEntity)
}
