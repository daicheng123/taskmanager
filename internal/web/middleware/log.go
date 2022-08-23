package middleware

import (
	"github.com/gin-gonic/gin"
	"taskmanager/pkg/logger"
	"time"
)

/*
	LoggerMiddleWare 自定义日志中间件
*/
type LoggerMiddleWare struct {
}

func NewLoggerMiddleWare() *LoggerMiddleWare {
	return &LoggerMiddleWare{}
}

func (lm *LoggerMiddleWare) OnRequest() gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now()
		clientIP := context.Request.RemoteAddr
		context.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		method := context.Request.Method
		reqUri := context.Request.RequestURI
		statusCode := context.Writer.Status()

		if statusCode >= 400 {
			logger.GinERROR(" %s %3d %13v %15s %s %s",
				startTime.Format("2006-01-02 15:04:05.9999"),
				statusCode,
				latencyTime,
				clientIP,
				method,
				reqUri)
		} else {
			logger.GinInfo(" %s %3d %13v %15s %s %s",
				startTime.Format("2006-01-02 15:04:05.9999"),
				statusCode,
				latencyTime,
				clientIP,
				method,
				reqUri)
		}
	}
}
