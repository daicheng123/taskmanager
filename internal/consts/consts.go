package consts

import "net/http"

const (
	UserControllerGroup = "user"
	AppManagerConfPath  = "APP_MANAGER_CONF_FILE"

	ManagerLog = iota // 服务日志
	GinLog            // gin框架日志
	TaskLog           // 任务日志
)

var StatusMessage = map[int]string{
	http.StatusOK:                  "请求成功",
	http.StatusUnauthorized:        "请求未认证",
	http.StatusBadRequest:          "请求参数错误",
	http.StatusInternalServerError: "服务内部错误",
}
