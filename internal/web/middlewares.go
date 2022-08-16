package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	webutils "taskmanager/internal/web/utils"
)

type MiddleWare interface {
	OnRequest(context *gin.Context) *webutils.Response
}

//CrossMiddleWare  跨域中间件
type CrossMiddleWare struct {
}

func NewCrossMiddleWare() *CrossMiddleWare {
	return &CrossMiddleWare{}
}

func (cross *CrossMiddleWare) OnRequest(context *gin.Context) *webutils.Response {
	rsp := new(webutils.Response)
	method := context.Request.Method
	if method != "" {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		context.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
	} else {
		rsp.Status = false
		rsp.Code = http.StatusMethodNotAllowed
		rsp.Message = "错误的请求方法"
		return rsp
	}

	if method == http.MethodOptions {
		context.AbortWithStatus(http.StatusNoContent)
	}
	rsp.Status = true
	return rsp
}
