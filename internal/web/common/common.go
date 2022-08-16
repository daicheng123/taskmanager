package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/service"
	"taskmanager/internal/web"
	webUtils "taskmanager/internal/web/utils"
)

type CommonController struct {
	userService *service.UserService
}

func NewCommonController() *CommonController {
	return &CommonController{
		userService: service.NewUserService(),
	}
}

func (cc *CommonController) CheckEmailExist(context *gin.Context) {
	email := context.Query("email")
	webUtils.Render(context, func() (interface{}, error) {
		return cc.userService.GetUserByEmail(email)
	}, http.StatusBadRequest)
}

func (cc *CommonController) GenEmailCode(context *gin.Context) {
	//email := context.Query("email")

	//webUtils.Render(context, func() (interface{}, error) {
	//	return cc.userService.GetUserByEmail(email)
	//}, http.StatusBadRequest)
}

func (cc *CommonController) Build(rc *web.RouterCenter) {
	commonGroup := rc.RG.Group("")
	commonGroup.Handle("GET", "/check_email_exists", cc.CheckEmailExist)
	commonGroup.Handle("POST", "/email_code", cc.CheckEmailExist)
}
