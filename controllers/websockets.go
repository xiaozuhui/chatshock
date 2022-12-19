package controllers

import (
	"chatshock/middlewares"
	"chatshock/services"
	"chatshock/utils"
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
	cr := engine.Group("/ws")
	cr.Use(middlewares.JWTAuth())
	cr.GET("", e.LinkWebsocket)
}

// LinkWebsocket 连接websocket
func (e *WebSocketController) LinkWebsocket(c *gin.Context) {
	userService := services.UserFactory()
	var UserID uuid.UUID
	if claims, ok := c.Get("claims"); ok {
		UserID = claims.(utils.UserClaims).UUID
	} else {
		panic(errors.WithStack(errors.New("用户不存在")))
	}
	userEntity, err := userService.UserRepo.FindUser(UserID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 获取到连接
	conn, err := websocket.Accept(c.Writer, c.Request,
		&websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(conn *websocket.Conn, code websocket.StatusCode, reason string) {
		err := conn.Close(code, reason)
		if err != nil {
			log.Error(err)
		}
	}(conn, websocket.StatusInternalError, "")
	// 获取超时十秒的子context，后续收发操作都使用该子context
	ctx, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()
	// 保存链接
	chatUser := websockets.NewUser(UserID, userEntity, conn)
	// 开启给用户发消息的goroutine
	go chatUser.SendMessage(ctx)
	// 阻塞获取ws数据
	err = chatUser.ReceiveMessage(ctx)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{})
}
