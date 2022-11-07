package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
)

type MailController struct {
}

func NewMailController() *MailController {
	return &MailController{}
}

func (mc *MailController) checkEmailExist(ctx *gin.Context) {
	email := ctx.Query("email")
	service := admin.MailService{}
	res := service.CheckMailExists(email)
	ctx.JSON(http.StatusOK, res)

}

func (mc *MailController) genEmailCode(ctx *gin.Context) {
	mailService := new(admin.MailService)
	if err := ctx.ShouldBindJSON(mailService); err == nil {
		res := mailService.GenMailCode(ctx)
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
	}
}

func (mc *MailController) Build(rc *web.RouterCenter) {
	commonGroup := rc.RG.Group("")
	commonGroup.Handle("GET", "/check_email_exists", mc.checkEmailExist)
	commonGroup.Handle("POST", "/email_code", mc.genEmailCode)
}
