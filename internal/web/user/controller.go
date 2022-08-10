package user

import (
	"github.com/gin-gonic/gin"
	"taskmanager/internal/consts"
	"taskmanager/internal/web"
)

type UserController struct {
	*web.Response
}

func NewUserController() *UserController {
	return &UserController{
		Response: new(web.Response),
	}
}

func (uc *UserController) UserLogin(context *gin.Context) {

}

func (uc *UserController) Build(rc *web.RouterCenter) {
	userGroup := rc.RG.Group(consts.UserControllerGroup)
	userGroup.Handle("GET", "/login", uc.UserLogin)
}
