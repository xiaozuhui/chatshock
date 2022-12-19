package controllers

import (
	"chatshock/middlewares"
	"github.com/gin-gonic/gin"
)

type ChatRoomController struct {
}

func (e *ChatRoomController) Router(engine *gin.Engine) {
	crH := engine.Group("/v1/chatroom")
	crH.Use(middlewares.JWTAuth())
	crH.GET(":room_id")
	crH.POST("")
	crH.DELETE("")
}
