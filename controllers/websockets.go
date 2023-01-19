package controllers

import (
	"chatshock/services"
	"chatshock/websockets"
	"context"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"nhooyr.io/websocket"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-12-15 10:19:31
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-15 10:44:58
 * @Description: 聊天室相关的websocket路由函数
 */

type WebSocketController struct {
}

func (e *WebSocketController) Router(engine *gin.Engine) {
	engine.GET("/ws", e.LinkWebsocket)
}

// LinkWebsocket 连接websocket
func (e *WebSocketController) LinkWebsocket(c *gin.Context) {
	userService := services.UserFactory()
	userID := c.DefaultQuery("user_id", "")
	userUUID, err := uuid.FromString(userID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userEntity, err := userService.UserRepo.FindUser(userUUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 获取到连接
	conn, err := websocket.Accept(c.Writer, c.Request,
		&websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		panic(errors.WithStack(err))
	}
	defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
		err := conn.Close(code, reason)
		if err != nil {
			log.Error(err)
		}
	}(conn, websocket.StatusInternalError, "")
	// 获取超时十秒的子context，后续收发操作都使用该子context
	ctx, cancel := context.WithTimeout(c, time.Second*60*60)
	defer cancel()
	// 保存链接
	chatUser := websockets.NewUser(userUUID, userEntity, conn)
	// 开启给用户发消息的goroutine
	go chatUser.SendMessage(ctx)
	// 发送当前在线用户列表
	userListMsg := websockets.OnlineUserListMessage()
	websockets.UserLock.Lock()
	for _, u := range websockets.BroadCaster.Users {
		u.MessageChannel <- userListMsg
	}
	websockets.UserLock.Unlock()
	// 阻塞获取ws数据
	err = chatUser.ReceiveMessage(ctx)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{})
}
