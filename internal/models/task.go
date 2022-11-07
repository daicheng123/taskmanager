package models

import (
	"fmt"
	"taskmanager/internal/consts"
	"taskmanager/internal/models/common"
)

// Task 任务执行表
type Task struct {
	ID            uint              `gorm:"primarykey" json:"id"`
	TaskName      string            `gorm:"type:string;size:256;uniqueIndex:idx_un_tn"`
	ScriptId      uint              `gorm:"type:uint;size:10"`
	Mode          consts.TaskMode   `gorm:"type:uint;size:1"`
	TaskStatus    consts.TaskStatus `gorm:"type:uint;size:1;default:6"`
	Operator      string            `gorm:"type:string;size:32"`
	Steps         []*TaskStep       `gorm:"foreignKey:TaskRefer; references:TaskName; constraint:OnDelete:CASCADE"`
	TaskStartTime common.CustomTime
	TaskEndTime   common.CustomTime
	TaskDelta     string `gorm:"type:string; size:16"`
}

func (t *Task) TableName() string {
	return "tasks"
}

func (t *Task) GenerateUniqKey() string {
	return fmt.Sprintf("%s", t.TaskName)
}

//TaskStep  任务执行步骤表
type TaskStep struct {
	ID         uint              `gorm:"primarykey" json:"id"`
	StepNumber uint              `gorm:"type:uint;size:3"`
	ExecutorID uint              `gorm:"type:uint"`
	StepStatus consts.TaskStatus `gorm:"type:uint;size:1"`
	//StepStdout    	string            		`gorm:"type:string; size:256"` // 关联缓存key
	//StepStderr    	string
	StepResultKey  string                `gorm:"type:string; size:128"` // 关联缓存key
	TransferStatus consts.TransferStatus `gorm:"type:uint;size:1;default:6"`
	StepStartTime  common.CustomTime
	StepEndTime    common.CustomTime
	StepDelta      string `gorm:"type:string; size:16"`
	TaskRefer      string
}

func (t *TaskStep) IsExecuting() bool { return t.StepStatus == consts.TaskExecuting }
func (t *TaskStep) IsSuccess() bool   { return t.StepStatus == consts.TaskSuccess }
func (t *TaskStep) IsFailed() bool    { return t.StepStatus == consts.TaskFailed }
func (t *TaskStep) IsTimeout() bool   { return t.StepStatus == consts.TaskTimeout }
func (t *TaskStep) IsAbort() bool     { return t.StepStatus == consts.TaskAbort }

func (t *TaskStep) TableName() string {
	return "task_steps"
}

func (t *TaskStep) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%d", t.ID, t.TaskRefer)
}
