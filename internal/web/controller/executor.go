package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/models"
	"taskmanager/internal/repo/mapper"
	"taskmanager/internal/schemas"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
)

const (
	ExecutorControllerGroup = "executors"
)

type ExecutorController struct {
	*admin.ExecutorService
}

func NewExecutorController() *ExecutorController {
	return &ExecutorController{
		ExecutorService: admin.NewExecutorService(),
	}
}

// executorAdd 新增执行器
func (ec *ExecutorController) executorAdd(ctx *gin.Context) {
	executorInfo := &schemas.ExecutorInfo{}
	err := ctx.ShouldBindJSON(executorInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, ec.ExecutorService.ExecutorAdd(ctx, executorInfo))
}

// distributeKey 分发管理主机密钥至执行器
func (ec *ExecutorController) distributeKey(ctx *gin.Context) {
	distribution := &schemas.DistributeReq{}
	err := ctx.ShouldBindJSON(distribution)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	filter := &models.Executor{
		BaseModel: models.BaseModel{ID: distribution.ID},
	}
	result, err := mapper.GetExecutorMapper().PreLoadFindOne(filter)
	if err != nil {
		logger.Error("查询executor失败,err:[%s]", err.Error())
		ctx.JSON(http.StatusOK, serializer.DBErr("获取executor失败", err))
		return
	}
	ctx.JSON(http.StatusOK, ec.ExecutorService.DistributeKey(ctx, result))
}

// executorLists
func (ec *ExecutorController) executorLists(ctx *gin.Context) {
	srv := &admin.ListService{}
	err := ctx.ShouldBindQuery(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	if srv.IsNotPage {
		_, executors, err := srv.ExecutorList()
		if err != nil {
			ctx.JSON(http.StatusOK, serializer.DBErr(err.Error(), err))
			return
		}
		ctx.JSON(http.StatusOK, &serializer.Response{Data: executors})
		return
	}
	ctx.JSON(http.StatusOK, srv.Lister(srv.ExecutorList))
}

func (ec *ExecutorController) executorDelete(ctx *gin.Context) {
	srv := &admin.ExecutorDelService{}
	err := ctx.ShouldBindUri(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.ExecutorDelete())
}

func (ec *ExecutorController) executorsBatchDelete(ctx *gin.Context) {
	var data = struct {
		IdList []uint `json:"idList" binding:"required"`
	}{}

	err := ctx.ShouldBindJSON(&data)
	if err == nil {
		srv := &admin.ExecutorDelService{}
		ctx.JSON(http.StatusOK, srv.ExecutorBatchDelete(data.IdList))
		return
	}
	ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
}

func (ec *ExecutorController) executorsRefresh(ctx *gin.Context) {
	// 默认状态不入库
	flag := ctx.DefaultQuery("flag", "create")
	srv := &admin.ExeTestService{}
	err := ctx.ShouldBindJSON(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.RefreshNode(flag))
}

func (ec *ExecutorController) executorUpdate(ctx *gin.Context) {
	executorInfo := &schemas.ExecutorInfo{}
	err := ctx.ShouldBindJSON(executorInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, ec.ExecutorService.ExecutorUpdate(ctx, executorInfo))
}

func (ec *ExecutorController) executorOption(ctx *gin.Context) {
	srv := &admin.ExecutorService{}
	ctx.JSON(http.StatusOK, srv.ExecutorOptions())
}

func (ec *ExecutorController) Build(rc *web.RouterCenter) {
	execGroup := rc.RG.Group(ExecutorControllerGroup)
	execGroup.Handle(http.MethodPost, "/add_executor", ec.executorAdd)
	execGroup.Handle(http.MethodPut, "/distribute_key", ec.distributeKey)
	execGroup.Handle(http.MethodGet, "/list_executor", ec.executorLists)
	execGroup.Handle(http.MethodDelete, "/del_executor/:id", ec.executorDelete)
	execGroup.Handle(http.MethodDelete, "/batch_delete", ec.executorsBatchDelete)
	execGroup.Handle(http.MethodPatch, "/refresh_status", ec.executorsRefresh)
	execGroup.Handle(http.MethodPut, "/update_executor", ec.executorUpdate)
	execGroup.Handle(http.MethodGet, "/option_executor", ec.executorOption)
}
