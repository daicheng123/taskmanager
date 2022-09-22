package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
)

const (
	ScriptControllerGroup = "scripts"
)

type ScriptController struct {
}

func NewScriptController() *ScriptController {
	return &ScriptController{}
}

func (sc *ScriptController) scriptAdd(ctx *gin.Context) {
	srv := &admin.ScriptService{}
	err := ctx.ShouldBind(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	srv.AddScript()
}

type script struct {
	Content string `json:"scriptContent" binding:"required"`
}

func (sc *ScriptController) scriptTest(ctx *gin.Context) {
	srv := &admin.ScriptService{}
	var s script
	err := ctx.ShouldBind(&s)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.CheckShellScript(s.Content))
}

func (sc *ScriptController) Build(rc *web.RouterCenter) {
	ScriptGroup := rc.RG.Group(ScriptControllerGroup)
	ScriptGroup.Handle(http.MethodPost, "/add_script", sc.scriptAdd)
	ScriptGroup.Handle(http.MethodPut, "/test_script", sc.scriptTest)
}
