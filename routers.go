/*
 * @Author: xiaozuhui
 * @Date: 2022-10-31 12:27:23
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-10-31 13:48:36
 * @Description:
 */
package main

import (
	"chatshock/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	new(controllers.UserController).Router(r)
	new(controllers.TokenController).Router(r)
}
