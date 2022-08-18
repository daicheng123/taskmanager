package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"taskmanager/internal/service"
	"taskmanager/internal/web"
	webUtils "taskmanager/internal/web/utils"
	"taskmanager/pkg/logger"
)

type CommonController struct {
	userService   *service.UserService
	commonService *service.CommonService
}

func NewCommonController() *CommonController {
	return &CommonController{
		userService:   service.NewUserService(),
		commonService: service.NewCommonService(),
	}
}

func (cc *CommonController) CheckEmailExist(context *gin.Context) {
	email := context.Query("email")
	webUtils.Render(context, func() (interface{}, error) {
		return cc.userService.GetUserByEmail(email)
	})
}

func (cc *CommonController) GenEmailCode(context *gin.Context) {
	data := struct {
		Email string `json:"email" binding:"required,email"`
	}{}

	webUtils.RenderNoBody(context, func() error {
		if err := context.ShouldBindJSON(&data); err != nil {
			logger.Error("email: %s", err.Error())
			return errors.New(webUtils.StatusMessage[webUtils.InvalidParamsCode])
		}
		return cc.commonService.GenEmailCode(data.Email)
	})
}

func (cc *CommonController) Build(rc *web.RouterCenter) {
	commonGroup := rc.RG.Group("")
	commonGroup.Handle("GET", "/check_email_exists", cc.CheckEmailExist)
	commonGroup.Handle("POST", "/email_code", cc.GenEmailCode)
}
