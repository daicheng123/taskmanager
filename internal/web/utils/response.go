package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Response 响应体
type Response struct {
	Code    int         `json:"code,omitempty"` // 兼容中间件code
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  bool        `json:"status,omitempty"`
}

//Render 返回值处理
func Render(ctx *gin.Context, fn func() (interface{}, error), codeSlice ...int) {
	var (
		result = &Response{
			Status:  true,
			Message: nil,
			Data:    nil,
			Code:    http.StatusInternalServerError,
		}
	)

	if codeSlice != nil && len(codeSlice) > 0 {
		result.Code = codeSlice[0]
	}

	defer func() {
		if p := recover(); p != nil {
			result.Status = false
			result.Message = fmt.Sprintf("%s", p)
		}
	}()

	data, bizErr := fn()
	if bizErr != nil {
		//result.Code = http.StatusOK
		result.Message = bizErr.Error()
		result.Status = false
		result.Data = data
		ctx.JSON(result.Code, result)
		return
	}

	result.Code = http.StatusOK
	result.Data = data
	ctx.JSON(result.Code, result)
	return
}
