package controllers

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 12:25:56
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 12:26:15
 * @Description:
 */

import (
	"chatshock/services"
	"chatshock/utils"
	"github.com/pkg/errors"
	"log"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
}

func (e *TokenController) Router(engine *gin.Engine) {
	token := engine.Group("/v1/token")
	token.POST("/refresh", e.RefreshToken) // 刷新Token
}

func (e *TokenController) RefreshToken(c *gin.Context) {
	tokenParam := &struct {
		RefreshTokenStr string `json:"refresh"`
	}{}
	err := c.Bind(tokenParam)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.JSON(403, gin.H{
			"code": -1,
			"msg":  "参数错误",
			"data": nil,
		})
		c.Abort()
	}
	claims, err := utils.ParseRefreshToken(tokenParam.RefreshTokenStr)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.JSON(500, gin.H{
			"code": -1,
			"msg":  err.Error(),
			"data": nil,
		})
		c.Abort()
	}
	token, refresh, expireTime, err := utils.GenerateToken(claims.PhoneNumber)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.JSON(500, gin.H{
			"code": -1,
			"msg":  err.Error(),
			"data": nil,
		})
		c.Abort()
	}
	token_ := services.MakeToken(token, refresh, *expireTime)
	c.JSON(200, gin.H{
		"code": -1,
		"msg":  "Token生成成功",
		"data": token_,
	})
}
