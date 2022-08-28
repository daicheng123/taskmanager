package serializer

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg"`
	Error   string      `json:"error,omitempty"`
}

// AppError 应用错误，实现了error接口
type AppError struct {
	Code     int
	Message  string
	RawError error
}

// NewError 返回新的错误对象
func NewError(code int, msg string, err error) AppError {
	return AppError{
		Code:     code,
		Message:  msg,
		RawError: err,
	}
}

// NewErrorFromResponse 从 serializer.Response 构建错误
func NewErrorFromResponse(resp *Response) AppError {
	return AppError{
		Code:     resp.Code,
		Message:  resp.Message,
		RawError: errors.New(resp.Error),
	}
}

// WithError 将应用error携带标准库中的error
func (err *AppError) WithError(raw error) AppError {
	err.RawError = raw
	return *err
}

// Error 返回业务代码确定的可读错误信息
func (err *AppError) Error() string {
	return err.Message
}

// ParamErr 各种参数错误
func ParamErr(msg string, err error) *Response {
	if msg == "" {
		msg = "参数错误"
	}
	return Err(CodeParamErr, msg, err)
}

// Err 通用错误处理
func Err(errCode int, msg string, err error) *Response {
	// 底层错误是AppError，则尝试从AppError中获取详细信息
	if appError, ok := err.(*AppError); ok {
		errCode = appError.Code
		err = appError.RawError
		msg = appError.Message
	}
	res := &Response{
		Code:    errCode,
		Message: msg,
	}
	// 生产环境隐藏底层报错

	if err != nil && gin.Mode() != gin.ReleaseMode {
		res.Error = err.Error()
	}
	return res
}

// DBErr 数据库操作失败
func DBErr(msg string, err error) *Response {
	if msg == "" {
		msg = "数据库操作失败"
	}

	return Err(CodeDBError, msg, err)
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	// CodeNotFullySuccess 未完全成功
	CodeNotFullySuccess = 203
	// CodeCheckLogin 未登录
	CodeCheckLogin = 401
	// CodeNoPermissionErr 未授权访问
	CodeNoPermissionErr = 403
	// CodeNotFound 资源未找到
	CodeNotFound = 404
	// CodeNotSet 未定错误，后续尝试从error中获取
	CodeNotSet = -1

	//CodeParamErr 各类参数错误
	CodeParamErr = 40001

	// CodeMailSendErr 邮件发送错误
	CodeMailSendErr = 4002

	// CodeServerInternalError 服务器内部错误
	CodeServerInternalError = 50001

	//CodeDBError  数据库操作失败
	CodeDBError = 50002
)
