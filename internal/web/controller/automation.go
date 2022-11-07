package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/schemas"
	"taskmanager/internal/service/automation"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
)

const (
	automationControllerGroup = "automation"
)

type AutomationController struct {
	operationTaskService *automation.OperationTaskService
}

func NewAutomationController() *AutomationController {
	return &AutomationController{
		operationTaskService: automation.NewOperationTaskService(),
	}
}

func (ac *AutomationController) executeOperation(ctx *gin.Context) {
	req := &schemas.OperationPayloadReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, ac.operationTaskService.ExecuteScript(ctx, req))
}

func (ac *AutomationController) fetchOperationTask(ctx *gin.Context) {
	req := &schemas.GetOperationTaskReq{}
	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, ac.operationTaskService.GetTaskDetail(ctx, req.TaskName))
}

func (ac *AutomationController) Build(rc *web.RouterCenter) {
	automationGroup := rc.RG.Group(automationControllerGroup)
	automationGroup.Handle("POST", "/execute_script", ac.executeOperation)
	automationGroup.Handle("PUT", "/fetch_operation_task/:taskName", ac.fetchOperationTask)
}
