package controllers

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 09:33:56
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 09:13:57
 * @Description:
 */

import (
	"chatshock/applications"
	"chatshock/entities"
	"chatshock/middlewares"
	"chatshock/services"
	"chatshock/utils"
	"log"

	"github.com/gofrs/uuid"

	"github.com/pkg/errors"

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

// Register
/**
 * @description: 注册用户
 * @param {*gin.Context} c
 * @return {*}
 * @author: xiaozuhui
 */
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
	userApplication := applications.NewUserApplication()
	userEntity := entities.UserEntity{
		NickName:    userParam.NickName,
		Password:    userParam.Password,
		PhoneNumber: userParam.PhoneNumber,
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
 * @return {*}
 * @author: xiaozuhui
 */
func (e *UserController) CheckValidCode(c *gin.Context) {
	userAuth := struct {
		PhoneNumber string `json:"phone_number"`
		ValidCode   string `json:"valid_code"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	err = utils.CheckValidCode(userAuth.PhoneNumber, userAuth.ValidCode)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, gin.H{"code": 1, "msg": "验证通过"})
}

// Login
/**
 * @description: 使用密码登录
 * @param {*gin.Context} c
 * @return {*}
 * @author: xiaozuhui
 */
func (e *UserController) Login(c *gin.Context) {
	userAuth := struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userService := services.UserFactory()
	// 1、检查账号密码
	isCheck, err := userService.CheckPassword(userAuth.PhoneNumber, userAuth.Password)
	if err != nil {
		panic(errors.WithStack(err))
	}
	if !isCheck {
		panic(errors.WithStack(errors.New("密码错误")))
	}
	// 2、登录，返回用户信息和token
	userInfo, err := userService.Login(userAuth.PhoneNumber)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, userInfo)
}

// LoginByPhoneNumber
/**
 * @description: 使用手机号登录
 * @param {*gin.Context} c
 * @return {*}
 * @author: xiaozuhui
 */
func (e *UserController) LoginByPhoneNumber(c *gin.Context) {
	userAuth := struct {
		PhoneNumber string `json:"phone_number"`
		ValidCode   string `json:"valid_code"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	err = utils.CheckValidCode(userAuth.PhoneNumber, userAuth.ValidCode)
	if err != nil {
		panic(errors.WithStack(err))
	}
	userService := services.UserFactory()
	userInfo, err := userService.Login(userAuth.PhoneNumber)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, userInfo)
}

// PhoneNumber
/**
 * @description: 手机号验证且发送验证码
 * @param {*gin.Context} c
 * @return {*}
 * @author: xiaozuhui
 */
func (e *UserController) PhoneNumber(c *gin.Context) {
	userAuth := struct {
		PhoneNumber string `json:"phone_number"`
	}{}
	err := c.Bind(&userAuth)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 检查该手机号是否已经注册
	userService := services.UserFactory()
	_, err = userService.GetUserByPhoneNumber(userAuth.PhoneNumber)
	if err != nil {
		panic(errors.WithStack(err))
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

// GetAccount
/**
 * @description: 获取账号信息
 * @param {*gin.Context} c
 * @return {*}
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
 * @return {*}
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

// UpdateAvatar 修改头像
func (e *UserController) UpdateAvatar(c *gin.Context) {
	userService := services.UserFactory()
	fileService := services.FileFactory()

	id := c.Param("id")
	UUID, err := uuid.FromString(id)
	if err != nil {
		panic(errors.WithStack(err))
	}
	avatar, err := c.FormFile("avatar")
	if err != nil {
		panic(errors.WithStack(err))
	}
	user, err := userService.GetUser(UUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	imgInfo, err := utils.UploadFiles(user.PhoneNumber, avatar.Filename, avatar)
	if err != nil {
		panic(errors.WithStack(err))
	}
	// 保存信息
	fileEntity, err := fileService.SaveFile(imgInfo, "photo", avatar.Header.Get("Content-Type"))
	if err != nil {
		return
	}
	userEntity := entities.UserEntity{
		BaseEntity: entities.BaseEntity{
			UUID: user.UUID,
		},
		PhoneNumber: user.PhoneNumber,
		Avatar:      fileEntity,
	}
	err = userService.UpdateAccount(&userEntity)
	if err != nil {
		panic(errors.WithStack(err))
	}
	user, err = userService.GetUser(UUID)
	if err != nil {
		panic(errors.WithStack(err))
	}
	c.JSON(200, user)
}

func (e *UserController) RebindPhoneNumber(c *gin.Context) {

}

func (e *UserController) ResetPassword(c *gin.Context) {

}
