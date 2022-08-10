package consts

import "net/http"

const (
	UserControllerGroup = "user"
)

var StatusMessage = map[int]string{
	http.StatusOK:                  "请求成功",
	http.StatusUnauthorized:        "请求未认证",
	http.StatusBadRequest:          "请求参数错误",
	http.StatusInternalServerError: "服务内部错误",
}
