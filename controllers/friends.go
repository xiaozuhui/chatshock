package controllers

import (
	"chatshock/applications"
	"chatshock/middlewares"
	"chatshock/services"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-08 21:47:27
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 08:23:44
 * @Description:
 */

type FriendController struct {
}

func (e *FriendController) Router(engine *gin.Engine) {
	friendGroup := engine.Group("/v1/user/friend")
	friendGroup.Use(middlewares.JWTAuth())
	friendGroup.GET("/:user_id", e.GetFriends)         // 获取该用户所有的好友
	friendGroup.POST("/:user_id", e.AddFriend)         // 发送好友申请
	friendGroup.DELETE("/:user_id", e.DelFriend)       // 删除好友
	friendGroup.POST("/:user_id/apply", e.ApplyFriend) // 同意好友申请
}

func (e *FriendController) GetFriends(c *gin.Context) {
	id := c.Param("user_id")
	UUID, err := uuid.FromString(id)
	if err != nil {
		panic(errors.WithStack(err))
	}
	application := applications.NewFriendApplication()
	friends, err := application.GetFriends(UUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, friends)
}

func (e *FriendController) AddFriend(c *gin.Context) {
	userID := c.Param("user_id")
	UserUUID, err := uuid.FromString(userID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	otherID, isExist := c.GetQuery("other_id")
	if !isExist {
		panic(errors.WithStack(errors.New("参数错误")))
	}
	OtherUUID, err := uuid.FromString(otherID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	friendService := services.FriendFactory()
	err = friendService.AddFriend(UserUUID, OtherUUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{})
}

func (e *FriendController) DelFriend(c *gin.Context) {
	userID := c.Param("user_id")
	UserUUID, err := uuid.FromString(userID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	otherID, isExist := c.GetQuery("other_id")
	if !isExist {
		panic(errors.WithStack(errors.New("参数错误")))
	}
	OtherUUID, err := uuid.FromString(otherID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	friendService := services.FriendFactory()
	err = friendService.DeleteFriend(UserUUID, OtherUUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{})
}

func (e *FriendController) ApplyFriend(c *gin.Context) {
	userID := c.Param("user_id")
	UserUUID, err := uuid.FromString(userID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	otherID, isExist := c.GetQuery("other_id")
	if !isExist {
		panic(errors.WithStack(errors.New("参数错误")))
	}
	OtherUUID, err := uuid.FromString(otherID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	friendService := services.FriendFactory()
	err = friendService.ApplyFriend(UserUUID, OtherUUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{})
}
