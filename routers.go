package main

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 12:27:23
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-15 10:41:50
 * @Description:
 */

import (
	"chatshock/controllers"
	"github.com/gin-gonic/gin"
)

// InitRouters 导入Controllers
func InitRouters(r *gin.Engine) {
	// http
	new(controllers.UserController).Router(r)
	new(controllers.TokenController).Router(r)
	new(controllers.FriendController).Router(r)
	// ws
	new(controllers.ChatRoomController).Router(r)
}
