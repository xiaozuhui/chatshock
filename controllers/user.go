package controllers

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:33:56
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 14:18:25
 * @Description:
 */

import (
	"chatshock/entities"
	"chatshock/middlewares"
	"chatshock/services"
	"chatshock/utils"
	"github.com/pkg/errors"
	"log"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (e *UserController) Router(engine *gin.Engine) {
	userGroup := engine.Group("/v1/user")
	userGroup.POST("/login", e.Login)                           // 可以使用用户名密码
	userGroup.POST("/login/phone_number", e.LoginByPhoneNumber) // 手机号登陆
	userGroup.POST("/register", e.Register)                     // 注册
	userGroup.POST("/valid/phone_number", e.PhoneNumber)        // 发送手机验证码
	userGroup.POST("/valid/check_valid_code", e.CheckValidCode) // 验证验证码

	accountGroup := engine.Group("/v1/account")
	accountGroup.Use(middlewares.JWTAuth())
	accountGroup.GET("/:id", e.GetAccount)                     // 获取账号信息
	accountGroup.PUT("/:id", e.UpdateAccount)                  // 更新账号信息
	accountGroup.POST("/:id/avatar", e.UpdateAvatar)           // 更新头像
	accountGroup.PUT("/:id/phone_number", e.RebindPhoneNumber) // 更新手机号
	accountGroup.PUT("/:id/reset_password", e.ResetPassword)   // 找回密码
}

// Register 注册用户
func (e *UserController) Register(c *gin.Context) {
	userParam := struct {
		NickName    string `json:"nickname"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
	}{}
	err := c.Bind(&userParam)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userService := services.UserFactory()
	userEntity := entities.UserEntity{
		NickName:    userParam.NickName,
		Password:    userParam.Password,
		PhoneNumber: userParam.PhoneNumber,
	}
	userInfo, err := userService.Register(userEntity)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, userInfo)
}

// CheckValidCode 检查验证码
func (e *UserController) CheckValidCode(c *gin.Context) {

}

func (e *UserController) Login(c *gin.Context) {

}

func (e *UserController) LoginByPhoneNumber(c *gin.Context) {

}

// PhoneNumber 手机号验证且发送验证码
func (e *UserController) PhoneNumber(c *gin.Context) {
	userAuth := struct {
		PhoneNumber string `json:"phone_number"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		c.JSON(403, gin.H{
			"code": -1,
			"msg":  "参数错误",
			"data": nil,
		})
		return
	}
	// 生成验证码
	validCode := utils.GenerateValidCode(utils.RegisterOrLogin)
	// 存入redis
	_, err = utils.RedisSet(userAuth.PhoneNumber, validCode.ValidCode, &validCode.ExpireTime)
	if err != nil {
		log.Fatalf("发生错误: %v", err)
	}
	// 去请求短信服务
	err = utils.SendValidMessage(userAuth.PhoneNumber, map[string]string{"code": validCode.ValidCode})
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 返回请求
	c.JSON(200, gin.H{"code": 1, "msg": "已发送验证码"})
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
