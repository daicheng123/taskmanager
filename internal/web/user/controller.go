package user

import (
	"github.com/gin-gonic/gin"
	"taskmanager/internal/consts"
	"taskmanager/internal/service"
	"taskmanager/internal/web"
	webUtils "taskmanager/internal/web/utils"
)

type UserController struct {
	*webUtils.Response
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		Response: new(webUtils.Response),
	}
}

func (uc *UserController) UserLogin(context *gin.Context) {

}

// UserRegister 用户注册地址
func (uc *UserController) UserRegister(context *gin.Context) {

}

func (uc *UserController) UserInfo(context *gin.Context) {

}

func (uc *UserController) Build(rc *web.RouterCenter) {
	userGroup := rc.RG.Group(consts.UserControllerGroup)
	userGroup.Handle("GET", "/login", uc.UserLogin)
	userGroup.Handle("POST", "/register", uc.UserRegister)
	userGroup.Handle("GET", "/userinfo", uc.UserInfo)
}
