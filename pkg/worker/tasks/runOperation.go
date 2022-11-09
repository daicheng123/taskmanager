package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hibiken/asynq"
	"sync"
	"taskmanager/internal/consts"
	"taskmanager/internal/models"
	"taskmanager/internal/models/common"
	"taskmanager/internal/repo/mapper"
	"taskmanager/internal/schemas"
	"taskmanager/pkg/ansible"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
	"time"
)

func RunAnsibleAdhoc() func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		logger.Info("start to execute task")
		var (
			buf     = task.Payload()
			payload = &schemas.OperationPayloadReq{}
		)

		err := json.Unmarshal(buf, payload)
		if err != nil {
			logger.Error("反序列化任务载荷失败, taskName: %s, err:[%s]", payload.TaskName, err.Error())
			return &EndpointError{Err: err, delay: DefaultDelay}
		}

		taskMapper := mapper.GetTaskMapper()
		taskObj := &models.Task{
			TaskName:      payload.TaskName,
			ScriptId:      payload.ScriptID,
			Mode:          consts.ManualTask,
			TaskStatus:    consts.TaskExecuting,
			Operator:      payload.TaskOperator,
			TaskStartTime: common.CustomTime(time.Now()),
		}

		// 获取所有执行器
		executorMapper := mapper.GetExecutorMapper()
		executors, err := executorMapper.FindWithRangeID(payload.Servers...)
		if err != nil {
			logger.Error("获取执行器失败, taskName: %s, err:[%s]", payload.TaskName, err.Error())
			return err
		}

		wl := len(*executors)
		var (
			taskSteps = make([]*models.TaskStep, 0, wl)
			ansiCfgs  = make([]*ansible.OperatorConfig, 0, wl)
		)

		for _, e := range *executors {
			taskStep := &models.TaskStep{
				ExecutorID:    e.ID,
				StepNumber:    1,
				StepStatus:    consts.TaskExecuting,
				StepStartTime: common.CustomTime(time.Now()),
				TaskRefer:     taskObj.TaskName,
			}

			ansiCfg := ansible.NewConfig(
				ansible.WithOvertime(payload.ScriptOverTime), ansible.WithTaskExecutor(e.ID),
				ansible.WithTaskName(payload.TaskName), ansible.WithTaskUser(e.Account.AccountName))

			ansiCfg.SetCallBack(func(taskName string) ansible.ResultCallback {
				op := &ansible.OperationCallbackPayload{
					TaskName: taskName,
					ScriptId: payload.ScriptID,
				}
				return ansible.NewOperationCallback(op)
			})

			ansiCfg.AddInventory(e.IPAddr)
			err = ansiCfg.RenderPlayBook(payload.ScriptContent, payload.ScriptTypeId)
			if err != nil {
				logger.Error("渲染operation 模板失败, taskName: %s, err:[%s]", payload.TaskName, err.Error())
				return err
			}

			ansiCfgs = append(ansiCfgs, ansiCfg)
			taskSteps = append(taskSteps, taskStep)
		}

		taskObj.Steps = taskSteps
		if err := taskMapper.CreateTask(taskObj); err != nil {
			logger.Error("创建任务失败, taskName: %s, err:[%s]", payload.TaskName, err.Error())
			return err
		}

		waitGroup := sync.WaitGroup{}
		for _, cfg := range ansiCfgs {
			waitGroup.Add(1)
			go func(config *ansible.OperatorConfig) {
				defer func() {
					waitGroup.Done()
				}()
				timeCtx, cancelFunc := context.WithTimeout(ctx, time.Second*time.Duration(payload.ScriptOverTime))
				defer cancelFunc()
				if err = ansible.NewOperationTask(config).Execute(timeCtx); err != nil {

					sf := &models.TaskStep{
						TaskRefer:  payload.TaskName,
						ExecutorID: config.Executor,
					}
					step, err := taskMapper.FindTaskStepOne(sf)
					if err != nil {
						logger.Error("获取任务 %s 步骤详情失败, 执行器 %d, err: [%s] ", payload.TaskName, config.Executor, err.Error())
						return
					}
					// 任务超时
					if errors.Is(err, ansible.OperationTimeoutError) {
						step.StepStatus = consts.TaskTimeout
					} else {
						logger.Error("任务步骤运行失败, taskName: %s, err:[%s]", payload.TaskName, err.Error())
						return
					}
				}
			}(cfg)
		}

		waitGroup.Wait()
		taskObj, err = taskMapper.PreLoadFindOne(ctx, taskObj)
		now := time.Now()
		taskObj.TaskDelta = utils.CalcTaskDelta(taskObj.TaskStartTime, now)
		taskObj.TaskEndTime = common.CustomTime(now)

		var (
			failures  = make([]uint, 0)
			successes = make([]uint, 0)
			//timeout   = make([]uint, 0)
			aborts = make([]uint, 0)
		)

		for _, step := range taskObj.Steps {
			if step.IsSuccess() {
				successes = append(successes, 1)
			}
			if step.IsAbort() {
				aborts = append(aborts, 1)
			}
			if step.IsFailed() {
				failures = append(failures, 1)
			}
		}
		if len(aborts) > 0 {
			taskObj.TaskStatus = consts.TaskAbort
		} else if len(failures) > 0 {
			taskObj.TaskStatus = consts.TaskFailed
			//} else if len(timeout) > 0 {
			//	taskObj.TaskStatus = consts.TaskAbort
		} else {
			taskObj.TaskStatus = consts.TaskSuccess
		}
		if err = taskMapper.UpdateTask(taskObj); err != nil {
			logger.Error("更新任务状态失败, taskName: %s, err:[%s]", payload.TaskName, err.Error())
			return err
		}
		return nil
	}
}
