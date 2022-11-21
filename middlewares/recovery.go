package middlewares

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				s := "panic recover err: %v"
				logger.Errorf(s, err)
				c.Abort()
			}
		}()
		c.Next()
	}
}
