package utils

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"taskmanager/pkg/serializer"
)

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
		"email":    "emial format error",
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
func ErrorResponse(err error) *serializer.Response {
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
		return serializer.ParamErr("json 序列化错误", err)
	}

	return serializer.ParamErr("参数错误", err)
}
