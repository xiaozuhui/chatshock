package controllers

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-19 15:19:34
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-02 02:11:12
 * @Description:
 */

import (
	"chatshock/applications"
	"chatshock/middlewares"
	"chatshock/services"
	"chatshock/utils"
	"chatshock/websockets"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type ChatRoomController struct {
}

func (e *ChatRoomController) Router(engine *gin.Engine) {
	crH := engine.Group("/v1/chatroom")
	crH.Use(middlewares.JWTAuth())
	crH.GET(":room_id")
	crH.POST("", e.CreateChatRoom)
	crH.DELETE("")
	crH.GET("/:room_id/add", e.AddCreateChatRoom)
}

// CreateChatRoom
/**
 * @description: 创建聊天室
 * @param {*gin.Context} c
 * @return {*}
 * @author: xiaozuhui
 */
func (e *ChatRoomController) CreateChatRoom(c *gin.Context) {
	// 获取User数据
	var userID uuid.UUID
	if claims, ok := c.Get("claims"); ok {
		userID = claims.(*utils.UserClaims).UUID
	} else {
		panic(errors.WithStack(errors.New("获取不到用户信息")))
	}
	chatRoomParam := struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}{}
	err := c.ShouldBind(&chatRoomParam)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 1、创建聊天室数据
	chatRoomApplication := applications.NewChatRoomApplication()
	roomEntity, err := chatRoomApplication.CreateChatRoom(userID, chatRoomParam.Name, chatRoomParam.Description)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 2、创建聊天室，加入broadcast
	chatRoom := websockets.NewChatRoom(roomEntity)
	chatRoom.ChatRoomLock.Lock()
	websockets.BroadCaster.ChatRooms[roomEntity.UUID] = &chatRoom
	chatRoom.ChatRoomLock.Unlock()
	// 3、聊天室监听
	go chatRoom.Listen()
	c.JSON(200, roomEntity)
}

// AddCreateChatRoom 加入聊天室
func (e *ChatRoomController) AddCreateChatRoom(c *gin.Context) {
	// 获取用户id
	var userID uuid.UUID
	if claims, ok := c.Get("claims"); ok {
		userID = claims.(*utils.UserClaims).UUID
	} else {
		panic(errors.WithStack(errors.New("获取不到用户信息")))
	}
	// 获取聊天室id
	id := c.Param("room_id")
	roomID, err := uuid.FromString(id)
	if err != nil {
		panic(errors.WithStack(err))
	}
	fmt.Println("判断是否在聊天室中...")
	// 判断是否在聊天室中...
	chatRoom, ok := websockets.BroadCaster.ChatRooms[roomID]
	if !ok {
		panic(errors.WithStack(errors.New(fmt.Sprintf("ID为[%v]的聊天室不存在", roomID))))
	}
	fmt.Println("判断是否在聊天室中...")
	ok = chatRoom.CanEnterChatRoom(userID)
	if !ok {
		panic(errors.WithStack(errors.New(fmt.Sprintf("用户[%v]已经在聊天室中", userID))))
	}
	fmt.Println("判断是否在聊天室中...")
	// 加入了聊天室
	chatRoomService := services.ChatRoomFactory()
	room, err := chatRoomService.IntoChatRoom(userID, roomID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 更新broadcast
	chatRoom.UpdateChatRoomUser(room)
	// 发送欢迎消息
	fmt.Println(len(chatRoom.EnteringChannel))
	chatRoom.EnteringChannel <- websockets.BroadCaster.Users[userID]
	c.JSON(200, gin.H{})
}
