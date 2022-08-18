package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StatusCode uint

const (
	SuccessCode             StatusCode = 20000
	InvalidParamsCode       StatusCode = 40001
	ServerInternalErrorCode StatusCode = 40002
	ServerUnknownERRORCode  StatusCode = 40003
	UnauthorizedCode        StatusCode = 40004
	PermissionDeniedCode    StatusCode = 40005
)

func (sc StatusCode) convHttpCode() int {
	switch sc {
	case InvalidParamsCode:
		return http.StatusBadRequest
	case ServerInternalErrorCode, ServerUnknownERRORCode:
		return http.StatusInternalServerError
	case UnauthorizedCode:
		return http.StatusUnauthorized
	case PermissionDeniedCode:
		return http.StatusForbidden
	default:
		return http.StatusOK
	}
}

var StatusMessage = map[StatusCode]string{
	SuccessCode:             "请求成功",
	InvalidParamsCode:       "请求参数错误",
	ServerInternalErrorCode: "服务内部错误",
	ServerUnknownERRORCode:  "未知错误",
	UnauthorizedCode:        "请求未认证",
}

var MessageStatus = map[string]StatusCode{
	"请求成功":   SuccessCode,
	"请求参数错误": InvalidParamsCode,
	"服务内部错误": ServerInternalErrorCode,
	"未知错误":   ServerUnknownERRORCode,
	"请求未认证":  UnauthorizedCode,
}

//Response 响应体
type Response struct {
	Code    StatusCode  `json:"code,omitempty"` // 兼容中间件code
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  bool        `json:"status,omitempty"`
}

//Render 返回值处理
func Render(ctx *gin.Context, fn func() (interface{}, error)) {
	var (
		result = &Response{
			Status:  true,
			Message: nil,
			Data:    nil,
			Code:    SuccessCode,
		}
	)

	defer func() {
		if p := recover(); p != nil {
			result.Status = false
			result.Message = fmt.Sprintf("%s", p)
			result.Code = ServerInternalErrorCode
			ctx.JSON(ServerInternalErrorCode.convHttpCode(), result)
			return
		}
	}()

	data, e := fn()
	if e != nil {
		errStr := e.Error()
		var httpCode int
		if c, ok := MessageStatus[errStr]; ok {
			result.Code = c
			httpCode = c.convHttpCode()
		} else {
			result.Code = ServerInternalErrorCode
			httpCode = result.Code.convHttpCode()
		}

		result.Message = errStr
		result.Status = false
		result.Data = data
		ctx.JSON(httpCode, result)
		return
	}

	result.Data = data
	//result.Message =
	ctx.JSON(http.StatusOK, result)
	return
}

//RenderNoBody 返回值处理
func RenderNoBody(ctx *gin.Context, fn func() error) {
	var (
		result = &Response{
			Status:  true,
			Message: nil,
			Data:    nil,
			Code:    SuccessCode,
		}
	)

	defer func() {
		if p := recover(); p != nil {
			result.Status = false
			result.Message = fmt.Sprintf("%s", p)
			result.Code = ServerInternalErrorCode
			ctx.JSON(ServerInternalErrorCode.convHttpCode(), result)
			return
		}
	}()

	e := fn()
	if e != nil {
		errStr := e.Error()
		var httpCode int
		if c, ok := MessageStatus[errStr]; ok {
			result.Code = c
			httpCode = c.convHttpCode()
		} else {
			result.Code = ServerInternalErrorCode
			httpCode = result.Code.convHttpCode()
		}

		result.Message = errStr
		result.Status = false
		ctx.JSON(httpCode, result)
		return
	}

	ctx.JSON(http.StatusOK, result)
	return
}
