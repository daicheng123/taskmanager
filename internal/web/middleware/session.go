package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/service"
	"taskmanager/internal/service/admin"
	"taskmanager/pkg/serializer"
)

/*
	SessionMiddleWare 自定义session中间件
*/

type SessionMiddleWare struct {
}

func NewSessionMiddleWare() *SessionMiddleWare {
	return &SessionMiddleWare{}
}

func (sm *SessionMiddleWare) OnRequest() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !AllowUri(context.Request.URL.Path) {
			token := context.Request.Header.Get("X-Token")
			if token == "" {
				context.AbortWithStatusJSON(http.StatusUnauthorized, serializer.Err(serializer.CodeCheckLogin, "请求未认证", nil))
				return
			}
			if service.SessionJudge(token) {
				context.Set(admin.UserTokenStr, token)
				context.Next()
			} else {
				context.AbortWithStatusJSON(http.StatusUnauthorized, serializer.Err(serializer.CodeCheckLogin, "请求未认证", nil))
				return
			}
		}
		context.Next()
	}
}
