package schemas

import (
	"taskmanager/internal/consts"
	"taskmanager/internal/models/common"
)

// OperationPayloadReq 脚本执行载荷
type OperationPayloadReq struct {
	Servers        []uint `json:"servers"`
	ScriptID       uint   `json:"scriptId"`
	ScriptContent  string
	ScriptTypeId   uint
	ScriptOverTime int
	TaskName       string `json:"taskName"`
	UniqueTag      string `json:"uniqueTag"`
	TaskOperator   string `json:"taskOperator"`
	TaskId         string
}

type GetOperationTaskReq struct {
	TaskName string `json:"taskName" uri:"taskName" binding:"required"`
}

func (op *OperationPayloadReq) ToBytes() []byte {
	return Struct2BytesSlice[*OperationPayloadReq](op)
}

type TaskStepInfo struct {
	ExecutorId   uint   `json:"executor_id"`
	ExecutorName string `json:"executor_name"`
	StepStdout   string `json:"step_stdout"`
}
type TaskInfoResp struct {
	Mode            consts.TaskMode   `json:"mode"`
	TaskStatus      consts.TaskStatus `json:"task_status"`
	Operator        string            `json:"operator"`
	TaskStartTime   common.CustomTime `json:"task_start_time"`
	TaskEndTime     common.CustomTime `json:"task_end_time"`
	TaskDelta       string            `json:"task_delta"`
	RunningServers  []*TaskStepInfo   `json:"running_servers"`
	SuccessServers  []*TaskStepInfo   `json:"success_servers"`
	FailedServers   []*TaskStepInfo   `json:"failed_servers"`
	OvertimeServers []*TaskStepInfo   `json:"overtime_servers"`
	AbortServers    []*TaskStepInfo   `json:"abort_servers"`
}

type DistributeReq struct {
	ID              uint   `json:"id" binding:"required,gte=1"`
	IPAddr          string `json:"ipAddr" binding:"required"`
	SSHPort         uint   `json:"sshPort" binding:"required,min=1,max=65535"`
	AccountName     string `json:"accountName" binding:"required"`
	AccountPassword string `json:"accountPassword" binding:"required"`
}
