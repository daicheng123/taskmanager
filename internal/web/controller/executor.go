package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/consts"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
)

type ExecutorController struct {
}

// executorAdd 新增执行器
func (ec *ExecutorController) executorAdd(ctx *gin.Context) {
	srv := &admin.ExecutorService{}
	err := ctx.ShouldBindJSON(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.ExecutorAdd())
}

func (ec *ExecutorController) Build(rc *web.RouterCenter) {
	exGroup := rc.RG.Group(consts.ExecutorControllerGroup)
	exGroup.Handle("POST", "/add", ec.executorAdd)
}
