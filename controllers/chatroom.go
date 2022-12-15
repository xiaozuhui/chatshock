package controllers

import (
	"chatshock/entities"
	"chatshock/middlewares"
	"chatshock/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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
	cr := engine.Group("/ws/chatroom")
	cr.Use(middlewares.JWTAuth())
	cr.GET(":room_id", e.EnterChat)
}

// EnterChat 进入聊天室
func (e *ChatRoomController) EnterChat(c *gin.Context) {
	var UserID uuid.UUID
	if claims, ok := c.Get("claims"); ok {
		UserID = claims.(utils.UserClaims).UUID
	} else {
		panic(errors.WithStack(errors.New("用户不存在")))
	}
	id := c.Param("room_id")
	roomID, err := uuid.FromString(id)
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
			
		}
	}(conn, websocket.StatusInternalError, "")
	// 保存链接
	chatUser := entities.ChatRoomUser{UserID: UserID, Conn: conn}
	// 获取数据或是推出数据
	for {
		var v interface{}
		wsjson.Read(c, conn, &v)
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
