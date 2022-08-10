package web

import (
	"github.com/gin-gonic/gin"
	"taskmanager/internal/consts"
)

//Controller  控制器接口
type Controller interface {
	Build(rc *RouterCenter)
}

//Response 响应体
type Response struct {
	Code    int         `json:"code,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  bool        `json:"status,omitempty"`
}

func (r *Response) Render(context *gin.Context, status bool, code int, data, message interface{}) {
	if message == nil {
		message = consts.StatusMessage[code]
	}
	context.JSON(code, gin.H{
		"status":  status,
		"data":    data,
		"message": message,
	})
}
