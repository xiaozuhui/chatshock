package controllers

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:33:56
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 14:18:25
 * @Description:
 */

import (
	"chatshock/middlewares"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (e *UserController) Router(engine *gin.Engine) {
	userGroup := engine.Group("/v1/user")
	userGroup.POST("/login", e.Login)                           // 可以使用用户名密码
	userGroup.POST("/login/phone_number", e.LoginByPhoneNumber) // 手机号登陆
	userGroup.POST("/register", e.Register)                     // 注册
	userGroup.POST("/phone_number", e.PhoneNumber)              // 发送手机验证码

	accountGroup := engine.Group("/v1/account")
	accountGroup.Use(middlewares.JWTAuth())
	accountGroup.GET("/:id", e.GetAccount)                     // 获取账号信息
	accountGroup.PUT("/:id", e.UpdateAccount)                  // 更新账号信息
	accountGroup.POST("/:id/avatar", e.UpdateAvatar)           // 更新头像
	accountGroup.PUT("/:id/phone_number", e.RebindPhoneNumber) // 更新手机号
	accountGroup.PUT("/:id/reset_password", e.ResetPassword)   // 找回密码
}

func (e *UserController) Register(c *gin.Context) {

}

func (e *UserController) Login(c *gin.Context) {

}

func (e *UserController) LoginByPhoneNumber(c *gin.Context) {

}

func (e *UserController) PhoneNumber(c *gin.Context) {

}

func (e *UserController) GetAccount(c *gin.Context) {

}

func (e *UserController) UpdateAccount(c *gin.Context) {

}

func (e *UserController) UpdateAvatar(c *gin.Context) {

}

func (e *UserController) RebindPhoneNumber(c *gin.Context) {

}

func (e *UserController) ResetPassword(c *gin.Context) {

}
