package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/pkg/logger"
	"time"
)

type MiddleWare interface {
	OnRequest(context *gin.Context)
}

/*
	CrossMiddleWare  跨域中间件
*/
type CrossMiddleWare struct {
}

func NewCrossMiddleWare() *CrossMiddleWare {
	return &CrossMiddleWare{}
}

func (cross *CrossMiddleWare) OnRequest(context *gin.Context) {
	method := context.Request.Method
	origin := context.Request.Header.Get("Origin")
	if origin != "" {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		context.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
	}
	if method == "OPTIONS" {
		context.AbortWithStatus(http.StatusNoContent)
	}
	context.Next()
}

/*
	LoggerMiddleWare 自定义日志中间件
*/
type LoggerMiddleWare struct {
}

func NewLoggerMiddleWare() *LoggerMiddleWare {
	return &LoggerMiddleWare{}
}

func (lm *LoggerMiddleWare) OnRequest(context *gin.Context) {
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

/*
	ErrorMiddleWare  自定义错误中间件
	统一处理 400状态码
*/
type ErrorMiddleWare struct {
}

func NewErrorMiddleWare() *ErrorMiddleWare {
	return &ErrorMiddleWare{}
}

func (em *ErrorMiddleWare) OnRequest(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if context.IsAborted() {
				context.Status(http.StatusOK)
			}
			switch errStr := err.(type) {
			case string:
				logger.GinERROR("remote: %s, uri: %s: error: [%s]",
					context.Request.RemoteAddr, context.Request.URL, errStr)
			}
		}
		context.AbortWithStatus(http.StatusBadRequest)
	}()
	context.Next()
}
