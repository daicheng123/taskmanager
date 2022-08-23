package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	CrossMiddleWare  跨域中间件
*/
type CrossMiddleWare struct {
}

func NewCrossMiddleWare() *CrossMiddleWare {
	return &CrossMiddleWare{}
}

func (cross *CrossMiddleWare) OnRequest() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")
		if origin != "" {
			context.Header("Access-Control-Allow-Origin", "*")
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			context.Header("Access-Control-Allow-Headers", "Origin,X-Token, X-Requested-With, Content-Type, Accept, Authorization")
			context.Header("Access-Control-Expose-Headers", "Content-Length, X-Token, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			context.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
