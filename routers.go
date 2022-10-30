package main

import (
	"chatshock/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	new(controllers.UserController).Router(r)
	new(controllers.TokenController).Router(r)
	new(controllers.VerifyController).Router(r)
	new(controllers.AccountController).Router(r)
}
