package controllers

import (
	"chatshock/middlewares"

	"github.com/gin-gonic/gin"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-08 21:47:27
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-08 22:11:53
 * @Description:
 */

type FriendController struct {
}

func (e *FriendController) Router(engine *gin.Engine) {
	friendGroup := engine.Group("/v1/user/friend")
	friendGroup.Use(middlewares.JWTAuth())
	friendGroup.GET("/:id", e.GetFriends)         // 获取该用户所有的好友
	friendGroup.POST("/:id/add", e.AddFriend)     // 发送好友申请
	friendGroup.POST("/:id/del", e.DelFriend)     // 删除好友
	friendGroup.POST("/:id/apply", e.ApplyFriend) // 同意好友申请
}

func (e *FriendController) GetFriends(c *gin.Context) {

}

func (e *FriendController) AddFriend(c *gin.Context) {

}

func (e *FriendController) DelFriend(c *gin.Context) {

}

func (e *FriendController) ApplyFriend(c *gin.Context) {

}
