package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
)

const (
	DangerCmdControllerGroup = "dangerous_command"
)

type DangerCmdController struct {
}

func NewDangerCmdController() *DangerCmdController {
	return &DangerCmdController{}
}

func (dcc *DangerCmdController) dangerCmdSave(ctx *gin.Context) {
	srv := &admin.DangerousCommandService{}
	err := ctx.ShouldBindJSON(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.DangerousCmdSave())
}

func (dcc *DangerCmdController) dangerCmdDelete(ctx *gin.Context) {
	srv := &admin.DangerousCommandDelService{}
	err := ctx.ShouldBindUri(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.DangerousCmdDelete())
}

func (dcc *DangerCmdController) dangerCmdList(ctx *gin.Context) {
	srv := &admin.ListService{}
	err := ctx.ShouldBind(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.DangerousCmdList())
}

func (dcc *DangerCmdController) Build(rc *web.RouterCenter) {
	dangerCmdGroup := rc.RG.Group(DangerCmdControllerGroup)
	dangerCmdGroup.Handle(http.MethodPost, "/add", dcc.dangerCmdSave)
	dangerCmdGroup.Handle(http.MethodPut, "/update", dcc.dangerCmdSave)
	dangerCmdGroup.Handle(http.MethodDelete, "/del/:id", dcc.dangerCmdDelete)
	dangerCmdGroup.Handle(http.MethodGet, "/list", dcc.dangerCmdList)
}
