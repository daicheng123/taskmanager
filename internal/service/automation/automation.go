package automation

import (
	"context"
	"encoding/json"
	"taskmanager/internal/cache"
	"taskmanager/internal/models"
	"taskmanager/internal/repo/mapper"
	"taskmanager/internal/schemas"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/pkg/worker"
	"taskmanager/pkg/worker/client"
	"taskmanager/pkg/worker/tasks"
	"taskmanager/utils"
)

type OperationTaskService struct {
	scriptRepo   *mapper.ScriptMapper
	taskRepo     *mapper.TaskMapper
	executorRepo *mapper.ExecutorMapper
}

func NewOperationTaskService() *OperationTaskService {
	return &OperationTaskService{
		scriptRepo:   mapper.GetScriptMapper(),
		taskRepo:     mapper.GetTaskMapper(),
		executorRepo: mapper.GetExecutorMapper(),
	}
}

func (oes *OperationTaskService) ExecuteScript(ctx context.Context, req *schemas.OperationPayloadReq) *serializer.Response {
	// 获取操作脚本
	scriptFilter := &models.Script{BaseModel: models.BaseModel{ID: req.ScriptID}}
	result, err := oes.scriptRepo.FindScriptOne(ctx, scriptFilter)
	if err != nil {
		logger.Error("查询脚本出错, err:[%s]", err.Error())
		return serializer.DBErr("查询脚本出错", err)
	}

	req.ScriptContent = result.ScriptContent
	req.ScriptTypeId = result.ScriptType
	req.ScriptOverTime = int(result.OverTime)

	buf, err := json.Marshal(req)
	if err != nil {
		return serializer.Err(serializer.CodeWriteCreateAsyncTaskError, "解析任务数据失败", err)
	}
	job := &client.Job{
		ID:       req.TaskName,
		Payload:  json.RawMessage(buf),
		Delay:    tasks.DefaultDelay,
		MaxRetry: 0,
	}
	//fmt.Printf("hahahahahahahaha%#v, Bytes: %s, error: %v\n", job.Payload, string(Bytes), err)

	err = client.DefaultQueue.Write(string(worker.RunAnsibleAdhoc), string(worker.AsyncQueue), job)
	if err != nil {
		logger.Error("任务 %s 入列失败, err:[%s]", worker.RunAnsibleAdhoc, err.Error())
		return serializer.Err(serializer.CodeWriteCreateAsyncTaskError, "创建任务失败", err)
	}
	return &serializer.Response{Message: "ok"}
}

func (oes *OperationTaskService) GetTaskDetail(ctx context.Context, taskName string) *serializer.Response {
	var (
		taskFilter   = &models.Task{TaskName: taskName}
		taskInfoResp = &schemas.TaskInfoResp{}
	)
	task, err := oes.taskRepo.PreLoadFindOne(ctx, taskFilter)
	if err != nil {
		return serializer.DBErr("获取任务信息失败", err)
	}
	taskInfoResp.TaskStatus = task.TaskStatus
	taskInfoResp.TaskDelta = task.TaskDelta
	taskInfoResp.TaskStartTime = task.TaskStartTime
	taskInfoResp.TaskEndTime = task.TaskEndTime
	taskInfoResp.Mode = task.Mode
	taskInfoResp.Operator = task.Operator

	c := cache.NewRedisCache(ctx)
	for _, step := range task.Steps {
		executorFilter := &models.Executor{
			BaseModel: models.BaseModel{
				ID: step.ExecutorID,
			},
		}
		executor, err := oes.executorRepo.FindOne(executorFilter)
		if err != nil {
			return serializer.DBErr("获取执行器信息失败", err)
		}

		stdOut, err := c.HGet(utils.BuilderStr(task.TaskName, ":", executor.IPAddr), "stdOut")
		if err != nil {
			return serializer.Err(serializer.CodeGetStepStdoutError, err.Error(), err)
		}
		stepInfo := &schemas.TaskStepInfo{
			ExecutorId:   step.ExecutorID,
			ExecutorName: executor.HostName,
			StepStdout:   stdOut,
		}
		if step.IsAbort() {
			taskInfoResp.AbortServers = append(taskInfoResp.AbortServers, stepInfo)
		}

		if step.IsSuccess() {
			taskInfoResp.SuccessServers = append(taskInfoResp.SuccessServers, stepInfo)
		}

		if step.IsFailed() {
			taskInfoResp.FailedServers = append(taskInfoResp.FailedServers, stepInfo)
		}

		if step.IsExecuting() {
			taskInfoResp.RunningServers = append(taskInfoResp.RunningServers, stepInfo)
		}

		if step.IsTimeout() {
			taskInfoResp.OvertimeServers = append(taskInfoResp.OvertimeServers, stepInfo)
		}
	}
	return &serializer.Response{Data: taskInfoResp, Message: "ok"}
}
