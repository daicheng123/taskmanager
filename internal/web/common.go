package web

import "github.com/gin-gonic/gin"

//Controller  控制器接口
type Controller interface {
	Build(rc *RouterCenter)
}

//Response 响应体
type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  bool
}

func (r *Response) Render(context *gin.Context, code int, data interface{}, status bool) {
	context.JSON(code, gin.H{
		"status":  status,
		"data":    data,
		"message": "",
	})
}
