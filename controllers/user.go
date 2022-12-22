package controllers

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:33:56
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-22 10:43:08
 * @Description:
 */

import (
	"chatshock/applications"
	"chatshock/entities"
	"chatshock/interfaces"
	"chatshock/middlewares"
	"chatshock/services"
	"chatshock/utils"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func (e *UserController) Router(engine *gin.Engine) {
	userGroup := engine.Group("/v1/user")
	userGroup.POST("/login", e.Login)                           // 可以使用用户名密码
	userGroup.POST("/login/by_contract", e.LoginBySender)       // 手机号登陆
	userGroup.POST("/register", e.Register)                     // 注册
	userGroup.POST("/valid/send_valid_code", e.SenderCheckCode) // 发送手机验证码
	userGroup.POST("/valid/check_valid_code", e.CheckValidCode) // 验证验证码

	accountGroup := engine.Group("/v1/account")
	accountGroup.Use(middlewares.JWTAuth())
	accountGroup.GET("/:id", e.GetAccount)                   // 获取账号信息
	accountGroup.PUT("/:id", e.UpdateAccount)                // 更新账号信息
	accountGroup.POST("/:id/avatar", e.UpdateAvatar)         // 更新头像
	accountGroup.PUT("/:id/reset_password", e.ResetPassword) // 找回密码
}

// Register
/**
 * @description: 注册用户
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *UserController) Register(c *gin.Context) {
	userParam := struct {
		NickName     string `json:"nickname"`
		Password     string `json:"password"`
		EmailAddress string `json:"email_address"`
	}{}
	err := c.Bind(&userParam)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 验证
	var sender interfaces.ISender
	if userParam.EmailAddress != "" {
		sender = utils.EmailAddress{EmailAddress: userParam.EmailAddress}
	} else {
		panic(errors.WithStack(errors.New("电子邮件不能为空")))
	}
	isCheck, err := utils.RedisStrGet(sender.String() + utils.EmailRegisterCode.String() + "_checked")
	if err != nil {
		panic(errors.WithStack(err))
	}
	if isCheck == nil || *isCheck != "OK" {
		panic(errors.WithStack(errors.New(fmt.Sprintf("手机号或是电子邮件 %s 验证错误", sender.String()))))
	}
	// 创建账号
	userApplication := applications.NewUserApplication()
	userID, err := uuid.NewV4()
	if err != nil {
		panic(errors.WithStack(err))
	}
	userEntity := entities.UserEntity{
		BaseEntity: entities.BaseEntity{
			UUID: userID,
		},
		NickName: userParam.NickName,
		Password: userParam.Password,
		Email:    userParam.EmailAddress,
	}
	userInfo, err := userApplication.Register(userEntity)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, userInfo)
}

// CheckValidCode
/**
 * @description: 检查验证码
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *UserController) CheckValidCode(c *gin.Context) {
	userAuth := struct {
		EmailAddress string `json:"email_address"`
		ValidCode    string `json:"valid_code"`
		CheckFor     string `json:"check_for"` // TODO 检查的类型
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	var sender interfaces.ISender
	if userAuth.EmailAddress != "" {
		sender = utils.EmailAddress{EmailAddress: userAuth.EmailAddress}
	} else {
		panic(errors.WithStack(errors.New("电子邮件不能为空")))
	}
	err = utils.CheckValidCode(sender, userAuth.ValidCode)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 将类型和电子邮件拼接后存入redis
	st := utils.ParseSendCode(userAuth.CheckFor)
	_, err = utils.RedisSet(sender.String()+st.String()+"_checked", "OK", nil)
	if err != nil {
		_ = utils.RedisDelete(sender.String() + "_checked")
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{"code": 1, "msg": "验证通过"})
}

// Login
/**
 * @description: 使用密码登录
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *UserController) Login(c *gin.Context) {
	userAuth := struct {
		PhoneNumber  string `json:"phone_number"`
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userService := services.UserFactory()
	// 1、检查账号密码
	var sender interfaces.ISender
	if userAuth.EmailAddress != "" {
		sender = utils.EmailAddress{EmailAddress: userAuth.EmailAddress}
	} else {
		panic(errors.WithStack(errors.New("手机号和电子邮件不能同时为空")))
	}
	isCheck, err := userService.CheckPassword(sender, userAuth.Password)
	if err != nil {
		panic(errors.WithStack(err))
	}
	if !isCheck {
		panic(errors.WithStack(errors.New("密码错误")))
	}
	// 2、登录，返回用户信息和token
	userInfo, err := userService.Login(sender)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, userInfo)
}

// LoginBySender
/**
 * @description: 使用手机号或电子邮件登录
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *UserController) LoginBySender(c *gin.Context) {
	userAuth := struct {
		EmailAddress string `json:"email_address"`
		ValidCode    string `json:"valid_code"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	var sender interfaces.ISender
	if userAuth.EmailAddress != "" {
		sender = utils.EmailAddress{EmailAddress: userAuth.EmailAddress}
	} else {
		panic(errors.WithStack(errors.New("电子邮件不能为空")))
	}
	err = utils.CheckValidCode(sender, userAuth.ValidCode)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userService := services.UserFactory()
	userInfo, err := userService.Login(sender)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, userInfo)
}

// SenderCheckCode
/**
 * @description: 手机号验证且发送验证码
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *UserController) SenderCheckCode(c *gin.Context) {
	userAuth := struct {
		EmailAddress string `json:"email_address"`
		SendFor      string `json:"send_for"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 检查该手机号是否已经注册
	userService := services.UserFactory()
	st := utils.ParseSendCode(userAuth.SendFor)
	var sender interfaces.ISender
	if userAuth.EmailAddress != "" {
		_, err := userService.GetUserByEmailAddress(userAuth.EmailAddress)
		if err != nil {
			panic(errors.WithStack(err))
		}
		sender = utils.EmailAddress{EmailAddress: userAuth.EmailAddress}
	} else {
		panic(errors.WithStack(errors.New("电子邮件不能为空")))
	}
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 生成验证码
	validCode := utils.GenerateValidCode(utils.RegisterOrLogin)
	// 存入redis
	err = utils.SetCheckValidCode(sender, validCode.ValidCode, validCode.ExpireTime)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 去请求短信服务
	err = sender.SendMessage(st.String(), "注册",
		utils.Options{
			FieldName:  func() string { return "code" },
			FieldValue: func() string { return validCode.ValidCode }})
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 返回请求
	c.JSON(200, gin.H{"code": 1, "msg": "已发送验证码"})
}

// GetAccount
/**
 * @description: 获取账号信息
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *UserController) GetAccount(c *gin.Context) {
	id := c.Param("id")
	UUID, err := uuid.FromString(id)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userService := services.UserFactory()
	user, err := userService.GetUser(UUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, user)
}

// UpdateAccount
/**
 * @description:  修改一些基本信息
                  目前可修改：
				  	  1、NickName
					  2、Introduction
					  3、Gender
 * @param {*gin.Context}
 * @author: xiaozuhui
*/
func (e *UserController) UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	UUID, err := uuid.FromString(id)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userParam := struct {
		NickName     string `json:"nickname"`     // 昵称
		Introduction string `json:"introduction"` // 介绍
		Gender       string `json:"gender"`       // 性别
	}{}
	err = c.Bind(&userParam)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userService := services.UserFactory()
	_, err = userService.GetUser(UUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userEntity := entities.UserEntity{
		BaseEntity: entities.BaseEntity{
			UUID: UUID,
		},
		NickName:     userParam.NickName,
		Introduction: userParam.Introduction,
		Gender:       entities.GenderTypeStr(userParam.Gender),
	}
	err = userService.UpdateAccount(&userEntity)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{})
}

// UpdateAvatar
/**
 * @description: 更新头像
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *UserController) UpdateAvatar(c *gin.Context) {
	userApp := applications.NewUserApplication()
	id := c.Param("id")
	UUID, err := uuid.FromString(id)
	if err != nil {
		panic(errors.WithStack(err))
	}
	avatar, err := c.FormFile("avatar")
	if err != nil {
		panic(errors.WithStack(err))
	}
	userResp, err := userApp.UpdateAvatar(UUID, avatar)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, userResp)
}

// TODO: 重新绑定邮箱
func (e *UserController) RebindEmail(c *gin.Context) {

}

// TODO: 重置密码
func (e *UserController) ResetPassword(c *gin.Context) {

}
