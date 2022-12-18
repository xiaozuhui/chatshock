package controllers

import (
	"chatshock/middlewares"
	"chatshock/services"
	"chatshock/utils"
	"chatshock/websockets"
	log "github.com/sirupsen/logrus"

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

type ChatRoomController struct {
}

func (e *ChatRoomController) Router(engine *gin.Engine) {
	cr := engine.Group("/ws")
	cr.Use(middlewares.JWTAuth())
	cr.GET("", e.LinkWebsocket)
	crH := engine.Group("/v1/chatroom")
	crH.Use(middlewares.JWTAuth())
	crH.GET(":room_id")
}

// LinkWebsocket 连接websocket
func (e *ChatRoomController) LinkWebsocket(c *gin.Context) {
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
	// 保存链接
	chatUser := websockets.NewUser(UserID, userEntity, conn)
	// 获取到用户加入的所有聊天室
	// 开启给用户发消息的goroutine
	go chatUser.SendMessage(c)
	// 获取数据或是推出数据
	for {

	}
	// ctx, cancel := context.WithTimeout(c, time.Second*10)
	// defer cancel()
	// var v interface{}
	// err = wsjson.Read(ctx, conn, &v)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// log.Printf("received: %v", v)
}
