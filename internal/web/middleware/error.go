package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/pkg/logger"
)

/*
	ErrorMiddleWare  自定义错误中间件
	统一处理 400状态码
*/
type ErrorMiddleWare struct {
}

func NewErrorMiddleWare() *ErrorMiddleWare {
	return &ErrorMiddleWare{}
}

func (em *ErrorMiddleWare) OnRequest() gin.HandlerFunc {
	return func(context *gin.Context) {
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
}
