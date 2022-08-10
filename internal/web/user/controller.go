package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/consts"
	"taskmanager/internal/service"
	"taskmanager/internal/web"
)

type UserController struct {
	*web.Response
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		Response: new(web.Response),
	}
}

func (uc *UserController) UserLogin(context *gin.Context) {
	//context.ShouldBindJSON()
	token := ""
	uc.Render(context, true, http.StatusOK, token, nil)
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
