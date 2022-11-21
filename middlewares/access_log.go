package middlewares

import (
	"bytes"
	"chatshock/configs"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"time"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyWriter := &AccessLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bodyWriter
		beginTime := time.Now()
		c.Next()
		endTime := time.Now()
		logger.WithFields(logger.Fields{
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWriter.body.String(),
		}).Debugf("access_log: method: %s, status_code: %d, begin_time: %v, end_time: %v, duration: %v",
			c.Request.Method,
			bodyWriter.Status(),
			beginTime.Format(configs.FormatTime),
			endTime.Format(configs.FormatTime),
			endTime.Sub(beginTime).Microseconds())
	}
}
