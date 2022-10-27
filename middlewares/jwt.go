package middlewares

/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 10:52:04
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 10:52:41
 * @Description:
 */

import (
	"chatshock/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		tokenStr, err := utils.GetToken(authHeader)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			ctx.Abort() //结束后续操作
		}
		//解析token包含的信息
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				ctx.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "token授权已过期，请重新申请授权",
					"data":   nil,
				})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "无效的Token",
			})
			ctx.Abort()
			return
		}
		// 将当前请求的claims信息保存到请求的上下文c上
		ctx.Set("claims", claims)
		ctx.Next() // 后续的处理函数可以用过ctx.Get("claims")来获取当前请求的用户信息
	}
}
