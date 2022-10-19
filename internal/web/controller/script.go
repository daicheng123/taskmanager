package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
	"time"
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
	ctx.JSON(http.StatusOK, srv.AddScript())
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

func (sc *ScriptController) scriptRetrieve(ctx *gin.Context) {
	var srv = &admin.RetrieveScriptService{}
	err := ctx.ShouldBindUri(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.RetrieveScript())
}

func (sc *ScriptController) scriptDelete(ctx *gin.Context) {
	var srv = &admin.RetrieveScriptService{}
	err := ctx.ShouldBindUri(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.DeleteScript())
}

func (sc *ScriptController) scriptAuditsList(ctx *gin.Context) {

}

func (sc *ScriptController) scriptUpdate(ctx *gin.Context) {
	var srv = &admin.ScriptService{}
	err := ctx.ShouldBindJSON(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.UpdateScript())
}

func (sc *ScriptController) scriptList(ctx *gin.Context) {
	srv := &admin.ListService{}
	err := ctx.ShouldBindQuery(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.Lister(srv.ScriptsList))
}

func (sc *ScriptController) scriptDebug(ctx *gin.Context) {
	srv := &admin.DebugScriptService{}
	err := ctx.ShouldBindJSON(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}

	timeoutCtx, cancel := context.WithTimeout(
		context.Background(), time.Second*time.Duration(int(srv.OverTime)))
	defer cancel()
	ctx.JSON(http.StatusOK, srv.Debug(timeoutCtx))
}

func (sc *ScriptController) Build(rc *web.RouterCenter) {
	ScriptGroup := rc.RG.Group(ScriptControllerGroup)
	ScriptGroup.Handle(http.MethodPost, "/add_script", sc.scriptAdd)
	ScriptGroup.Handle(http.MethodPut, "/test_script", sc.scriptTest)
	ScriptGroup.Handle(http.MethodGet, "/list_script", sc.scriptList)
	ScriptGroup.Handle(http.MethodGet, "/retrieve_script/:id", sc.scriptRetrieve)
	ScriptGroup.Handle(http.MethodDelete, "/delete_script/:id", sc.scriptDelete)
	ScriptGroup.Handle(http.MethodPut, "/update_script", sc.scriptUpdate)
	ScriptGroup.Handle(http.MethodPost, "/debug_script", sc.scriptDebug)
}
