package controllers

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 12:25:56
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-12-06 09:42:38
 * @Description:
 */

import (
	"chatshock/services/resp"
	"chatshock/utils"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
}

func (e *TokenController) Router(engine *gin.Engine) {
	token := engine.Group("/v1/token")
	token.POST("/refresh", e.RefreshToken) // 刷新Token
}

// RefreshToken
/**
 * @description: 刷新Token，使用refresh token来刷新新的token
 * @param {*gin.Context} c
 * @author: xiaozuhui
 */
func (e *TokenController) RefreshToken(c *gin.Context) {
	tokenParam := &struct {
		RefreshTokenStr string `json:"refresh"`
	}{}
	err := c.Bind(tokenParam)
	if err != nil {
		panic(errors.WithStack(err))
	}
	claims, err := utils.ParseRefreshToken(tokenParam.RefreshTokenStr)
	if err != nil {
		panic(errors.WithStack(err))
	}
	token, refresh, expireTime, err := utils.GenerateToken(claims.PhoneNumber)
	if err != nil {
		panic(errors.WithStack(err))
	}
	tokenRes := resp.MakeToken(token, refresh, *expireTime)
	c.JSON(200, gin.H{
		"code": -1,
		"msg":  "Token生成成功",
		"data": tokenRes,
	})
}
