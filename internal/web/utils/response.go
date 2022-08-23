package utils

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"taskmanager/utils/serializer"
)

type StatusCode uint

//const (
//	SuccessCode             StatusCode = 20000
//	InvalidParamsCode       StatusCode = 40001
//	ServerInternalErrorCode StatusCode = 40002
//	ServerUnknownERRORCode  StatusCode = 40003
//	UnauthorizedCode        StatusCode = 40004
//	PermissionDeniedCode    StatusCode = 40005
//)
//
//func (sc StatusCode) convHttpCode() int {
//	switch sc {
//	case InvalidParamsCode:
//		return http.StatusBadRequest
//	case ServerInternalErrorCode, ServerUnknownERRORCode:
//		return http.StatusInternalServerError
//	case UnauthorizedCode:
//		return http.StatusUnauthorized
//	case PermissionDeniedCode:
//		return http.StatusForbidden
//	default:
//		return http.StatusOK
//	}
//}
//
//var StatusMessage = map[StatusCode]string{
//	SuccessCode:             "请求成功",
//	InvalidParamsCode:       "请求参数错误",
//	ServerInternalErrorCode: "服务内部错误",
//	ServerUnknownERRORCode:  "未知错误",
//	UnauthorizedCode:        "请求未认证",
//}
//
//var MessageStatus = map[string]StatusCode{
//	"请求成功":   SuccessCode,
//	"请求参数错误": InvalidParamsCode,
//	"服务内部错误": ServerInternalErrorCode,
//	"未知错误":   ServerUnknownERRORCode,
//	"请求未认证":  UnauthorizedCode,
//}
//
////Response 响应体
//type Response struct {
//	Code    StatusCode  `json:"code,omitempty"` // 兼容中间件code
//	Message interface{} `json:"message,omitempty"`
//	Data    interface{} `json:"data,omitempty"`
//	Status  bool        `json:"status,omitempty"`
//}
//
////Render 返回值处理
//func Render(ctx *gin.Context, fn func() (interface{}, error)) {
//	var (
//		result = &Response{
//			Status:  true,
//			Message: nil,
//			Data:    nil,
//			Code:    SuccessCode,
//		}
//	)
//
//	defer func() {
//		if p := recover(); p != nil {
//			result.Status = false
//			result.Message = fmt.Sprintf("%s", p)
//			result.Code = ServerInternalErrorCode
//			ctx.JSON(ServerInternalErrorCode.convHttpCode(), result)
//			return
//		}
//	}()
//
//	data, e := fn()
//	if e != nil {
//		errStr := e.Error()
//		var httpCode int
//		if c, ok := MessageStatus[errStr]; ok {
//			result.Code = c
//			httpCode = c.convHttpCode()
//		} else {
//			result.Code = ServerInternalErrorCode
//			httpCode = result.Code.convHttpCode()
//		}
//
//		result.Message = errStr
//		result.Status = false
//		result.Data = data
//		ctx.JSON(httpCode, result)
//		return
//	}
//
//	result.Data = data
//	//result.Message =
//	ctx.JSON(http.StatusOK, result)
//	return
//}
//
////RenderNoBody 返回值处理
//func RenderNoBody(ctx *gin.Context, fn func() error) {
//	var (
//		result = &Response{
//			Status:  true,
//			Message: nil,
//			Data:    nil,
//			Code:    SuccessCode,
//		}
//	)
//
//	defer func() {
//		if p := recover(); p != nil {
//			result.Status = false
//			result.Message = fmt.Sprintf("%s", p)
//			result.Code = ServerInternalErrorCode
//			ctx.JSON(ServerInternalErrorCode.convHttpCode(), result)
//			return
//		}
//	}()
//
//	e := fn()
//	if e != nil {
//		errStr := e.Error()
//		var httpCode int
//		if c, ok := MessageStatus[errStr]; ok {
//			result.Code = c
//			httpCode = c.convHttpCode()
//		} else {
//			result.Code = ServerInternalErrorCode
//			//httpCode = result.Code.convHttpCode()
//			httpCode = http.StatusOK
//		}
//
//		result.Message = errStr
//		result.Status = false
//		ctx.JSON(httpCode, result)
//		return
//	}
//
//	ctx.JSON(http.StatusOK, result)
//	return
//}

// ValidatorErrorMsg 根据Validator返回的错误信息给出错误提示
func ValidatorErrorMsg(filed string, tag string) string {
	// 未通过验证的表单域与中文对应
	fieldMap := map[string]string{
		"Email":    "Email",
		"Password": "Password",
		//"Path":     "Path",
		//"SourceID": "Source resource",
		//"URL":      "URL",
		//"Nick":     "Nickname",
	}
	// 未通过的规则与中文对应
	tagMap := map[string]string{
		"required": "cannot be empty",
		"min":      "too short",
		"max":      "too long",
		"email":    "format error",
	}
	fieldVal, findField := fieldMap[filed]
	if !findField {
		fieldVal = filed
	}
	tagVal, findTag := tagMap[tag]
	if findTag {
		// 返回拼接出来的错误信息
		return fieldVal + " " + tagVal
	}
	return ""
}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) serializer.Response {
	// 处理 Validator 产生的错误
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			return serializer.ParamErr(
				ValidatorErrorMsg(e.Field(), e.Tag()),
				err,
			)
		}
	}

	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.ParamErr("JSON marshall error", err)
	}

	return serializer.ParamErr("Parameter error", err)
}
