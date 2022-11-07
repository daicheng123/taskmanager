package ansible

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"taskmanager/internal/cache"
	"taskmanager/internal/consts"
	"taskmanager/internal/models/common"
	"taskmanager/internal/repo/mapper"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
	"time"
)

type ResultCallback interface {
	ParseResult(context.Context, *bytes.Buffer)
	RunnerOnSuccess(context.Context, string, string, *AnsiblePlaybookJSONResultsPlayTaskHostsItem)
	RunnerOnUnReachable(context.Context, string, string, *AnsiblePlaybookJSONResultsPlayTaskHostsItem)
	RunnerOnFailed(context.Context, string, string, *AnsiblePlaybookJSONResultsPlayTaskHostsItem)
}

type Result interface {
	map[string]interface{}
}

type CallbackBase struct {
	ResultCallback
	Results *AnsibleJSONResults
	lock    sync.Mutex
}

func NewBaseCallback() *CallbackBase {
	cb := &CallbackBase{
		Results: &AnsibleJSONResults{},
	}
	return cb
}

func (cb *CallbackBase) CheckTaskHostsItemResults(info *AnsiblePlaybookJSONResultsPlayTaskHostsItem) bool {
	if info.AnsiblePlaybookJSONResultsPlayTaskHostsItemResults != nil {
		return true
	}
	return false
}

func (cb *CallbackBase) ParseResult(ctx context.Context, buff *bytes.Buffer) {
	var err error
	cb.Results, err = JSONParse(buff.Bytes())
	if err != nil {
		logger.Error("解析ansible 任务执行结果错误, err:[%s]", err.Error())
		return
	}

	for _, play := range cb.Results.Plays {
		for _, task := range play.Tasks {
			//if task.Task.Name == consts.ExecuteScript || task.Task.Name == consts.TransferScript{
			for hostAddr, hostInfo := range task.Hosts {
				if cb.CheckTaskHostsItemResults(hostInfo) {
					if hostInfo.Results[0].Failed {
						cb.RunnerOnFailed(ctx, task.Task.Name, hostAddr, hostInfo)
					} else if hostInfo.Results[0].Unreachable {
						cb.RunnerOnUnReachable(ctx, task.Task.Name, hostAddr, hostInfo)
					} else {
						cb.RunnerOnSuccess(ctx, task.Task.Name, hostAddr, hostInfo)
					}
				} else {
					cb.RunnerOnSuccess(ctx, task.Task.Name, hostAddr, hostInfo)
				}
			}
			//}
		}
	}
	return
}

type OperationCallbackPayload struct {
	TaskName string
	ScriptId uint
}

type OperationCallback struct {
	*OperationCallbackPayload
	*CallbackBase
}

func NewOperationCallback(payload *OperationCallbackPayload) *OperationCallback {
	var (
		operationCallback = new(OperationCallback)
	)
	operationCallback.OperationCallbackPayload = payload

	baseCallBack := NewBaseCallback()
	baseCallBack.ResultCallback = operationCallback
	operationCallback.CallbackBase = baseCallBack
	return operationCallback
}

func (ac *OperationCallback) RunnerOnSuccess(ctx context.Context, ansibleTask string, hostAddr string, info *AnsiblePlaybookJSONResultsPlayTaskHostsItem) {
	fmt.Println(hostAddr, info, "Success")
	if ansibleTask != consts.DeleteScript {
		step, err := mapper.GetTaskMapper().FindStepByExecutorAndTask(hostAddr, ac.TaskName)
		if err != nil {
			logger.Error("获取执行步骤失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
			return
		}

		baseKey := utils.BuilderStr(ac.TaskName, ":", hostAddr)
		step.StepResultKey = baseKey

		if ansibleTask == consts.TransferScript {
			step.TransferStatus = consts.TransferSuccess
			info.Stdout = "文件传输成功"
		} else {
			step.StepStatus = consts.TaskSuccess
			now := time.Now()
			step.StepEndTime = common.CustomTime(now)
			step.StepDelta = utils.CalcTaskDelta(step.StepStartTime, now)
		}
		go utils.RunSafeWithMsg(func() {
			rc := cache.NewRedisCache(ctx)
			_, err = rc.HSet(baseKey, "stdOut", info.Stdout)
			//_, err = rc.HSet(baseKey, "stdErr", info.Stderr)
		}, "")

		if err = mapper.GetTaskMapper().UpdateTaskStep(step); err != nil {
			logger.Error("更新执行步骤状态失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
			return
		}
	}
}

func (ac *OperationCallback) RunnerOnUnReachable(ctx context.Context, ansibleTask string, hostAddr string, info *AnsiblePlaybookJSONResultsPlayTaskHostsItem) {
	//if ansibleTask != consts.DeleteScript {
	//	if ac.CheckTaskHostsItemResults(info) {
	//		info.Stdout = info.Results[0].Message
	//	}
	//
	//	step, err := mapper.GetTaskMapper().FindStepByExecutorAndTask(hostAddr, ac.TaskName)
	//	if err != nil {
	//		logger.Error("获取执行步骤失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
	//		return
	//	}
	//	step.StepStatus = consts.TaskFailed
	//
	//	if ansibleTask == consts.TransferScript {
	//		step.TransferStatus = consts.TransferFailed
	//	}
	//
	//	baseKey := utils.BuilderStr(ac.TaskName, ":", hostAddr)
	//	step.StepResultKey = baseKey
	//	// 解析结果存入redis
	//	go utils.RunSafeWithMsg(func() {
	//		rc := cache.NewRedisCache(ctx)
	//		_, err = rc.HSet(baseKey, "stdOut", info.Stdout)
	//		_, err = rc.HSet(baseKey, "stdErr", info.Stderr)
	//	}, "")
	//
	//	if err = mapper.GetTaskMapper().UpdateTaskStep(step); err != nil {
	//		logger.Error("更新执行步骤状态失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
	//		return
	//	}
	//}
	ac.onFailed(ctx, ansibleTask, hostAddr, info)
}

func (ac *OperationCallback) RunnerOnFailed(ctx context.Context, ansibleTask string, hostAddr string, info *AnsiblePlaybookJSONResultsPlayTaskHostsItem) {

	//if ansibleTask != consts.DeleteScript {
	//	fmt.Println(hostAddr, info, "Failed")
	//	if ac.CheckTaskHostsItemResults(info) {
	//		info.Stdout = info.Results[0].Message
	//	}
	//	step, err := mapper.GetTaskMapper().FindStepByExecutorAndTask(hostAddr, ac.TaskName)
	//	if err != nil {
	//		logger.Error("获取执行步骤失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
	//		return
	//	}
	//	baseKey := utils.BuilderStr(ac.TaskName, ":", hostAddr)
	//	step.StepResultKey = baseKey
	//	step.StepStatus = consts.TaskFailed
	//
	//	if ansibleTask == consts.TransferScript {
	//		step.TransferStatus = consts.TransferFailed
	//	}
	//	go utils.RunSafeWithMsg(func() {
	//		rc := cache.NewRedisCache(ctx)
	//		_, err = rc.HSet(baseKey, "stdOut", info.Stdout)
	//		_, err = rc.HSet(baseKey, "stdErr", info.Stderr)
	//	}, "")
	//
	//	if err = mapper.GetTaskMapper().UpdateTaskStep(step); err != nil {
	//		logger.Error("更新执行步骤状态失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
	//		return
	//	}
	//}
	ac.onFailed(ctx, ansibleTask, hostAddr, info)
}

func (ac *OperationCallback) onFailed(ctx context.Context, ansibleTask string, hostAddr string, info *AnsiblePlaybookJSONResultsPlayTaskHostsItem) {
	if ansibleTask != consts.DeleteScript {
		if ac.CheckTaskHostsItemResults(info) {
			info.Stdout = info.Results[0].Message
		}

		step, err := mapper.GetTaskMapper().FindStepByExecutorAndTask(hostAddr, ac.TaskName)
		if err != nil {
			logger.Error("获取执行步骤失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
			return
		}
		step.StepStatus = consts.TaskFailed

		if ansibleTask == consts.TransferScript {
			step.TransferStatus = consts.TransferFailed
		}

		baseKey := utils.BuilderStr(ac.TaskName, ":", hostAddr)
		step.StepResultKey = baseKey
		// 解析结果存入redis
		go utils.RunSafeWithMsg(func() {
			rc := cache.NewRedisCache(ctx)
			_, err = rc.HSet(baseKey, "stdOut", info.Stdout)
			_, err = rc.HSet(baseKey, "stdErr", info.Stderr)
		}, "")

		if err = mapper.GetTaskMapper().UpdateTaskStep(step); err != nil {
			logger.Error("更新执行步骤状态失败, 任务: %s, ansible 任务：%s err: [%s]", ac.TaskName, ansibleTask, err.Error())
			return
		}
	}
}
