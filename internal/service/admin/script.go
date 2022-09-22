package admin

import (
	"os"
	"taskmanager/internal/dal/mapper"
	"taskmanager/internal/models"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
)

type ScriptService struct {
	ScriptName    string `json:"scriptName" binding:"required"`
	ScriptContent string `json:"scriptContent" binding:"required"`
	ScriptType    uint   `json:"scriptType" binding:"required"`
	OverTime      uint   `json:"overTime" binding:"required"`
	ScriptTag     uint   `json:"scriptTag" binding:"required"`
	LastOperator  string `json:"lastOperator"`
	Remarks       string `json:"remarks"`
}

func (ss *ScriptService) AddScript() *serializer.Response {
	// 判断脚本是否存在危险命令
	cmdFilter := &models.DangerousCmd{}
	commands, err := mapper.GetDangerCmdMapper().ListAllDangerousCommand(cmdFilter)
	if err != nil {
		return serializer.DBErr("查询危险命令出错", err)
	}
	/*
		如果包含危险命令==> 进入审核(脚本状态：审核中)
		未包含危险命令===> 直接审核通过(脚本状态：通过)
	*/
	var needAudit bool
	for _, c := range commands {
		if c.CheckExistInScript(ss.ScriptContent) {
			needAudit = true
		}
	}

	return &serializer.Response{}
}

// CheckShellScript 校验shell类型脚本
func (ss *ScriptService) CheckShellScript(script string) *serializer.Response {
	tempName := "/tmp/" + utils.NewUuid() + ".sh"
	file, err := os.OpenFile(tempName, os.O_WRONLY|os.O_CREATE, os.FileMode.Perm(0644))
	defer file.Close()
	if err != nil {
		return serializer.Err(serializer.CodeServerInternalError, "打开测试脚本文件失败", err)
	}
	_, err = file.WriteString(script)
	if err != nil {
		return serializer.Err(serializer.CodeServerInternalError, "写入测试脚本文件失败", err)
	}
	out, err := utils.ExecuteShell("shellcheck " + tempName)
	if err != nil {
		logger.Error("执行脚本 %s 出错, err:%v", tempName, err)
	}
	return &serializer.Response{Data: out}
}
