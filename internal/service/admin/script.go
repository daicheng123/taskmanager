package admin

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"taskmanager/internal/consts"
	"taskmanager/internal/models"
	"taskmanager/internal/repo/mapper"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
)

type ScriptService struct {
	*models.Script `json:",inline"  binding:"required"`
	UserId         uint `json:"userId" binding:"omitempty,gte=0"`
}

type RetrieveScriptService struct {
	ID uint `uri:"id" binding:"required"`
}

type DebugScriptService struct {
	ServerId      uint   `json:"serverId" binding:"required"`
	OverTime      uint   `json:"overTime"  binding:"required"`
	ScriptContent string `json:"scriptContent" binding:"required"`
	ScriptType    uint   `json:"scriptType" binding:"required"`
}

func (ds *DebugScriptService) Debug(ctx context.Context) *serializer.Response {
	err, needAudit := CheckDangerCmd(ds.ScriptContent)
	if err != nil {
		return serializer.Err(serializer.CodeDangerCmdQueryError, err.Error(), err)
	}
	if needAudit {
		return serializer.Err(serializer.CodeDangerCmdDebugError, "脚本包含危险操作符,无法调试", nil)
	}

	var (
		srcPath, command string
	)
	filter := &models.Executor{
		BaseModel: models.BaseModel{
			ID: ds.ServerId,
		},
	}
	executor, err := mapper.GetExecutorMapper().PreLoadFindOne(filter)
	if err != nil {
		return serializer.DBErr("查询执行器出错", err)
	}
	password, err := utils.Decrypt(executor.Account.AccountPassword)
	if err != nil {
		return serializer.Err(serializer.CodeHostPasswordDecodeErr, err.Error(), err)
	}
	uuid := utils.NewUuid()
	if ds.ScriptType == consts.Python {
		srcPath = utils.BuilderStr("/tmp/", uuid, ".py")
		command = "python3 "
	} else if ds.ScriptType == consts.Shell {
		srcPath = utils.BuilderStr("/tmp/", uuid, ".sh")
		command = "bash "
	}
	dstPath := filepath.Join(executor.ExecutePath, filepath.Base(srcPath))
	command = command + dstPath
	// 写文件
	err = utils.WriteFile(srcPath, ds.ScriptContent)
	if err != nil {
		return serializer.Err(serializer.CodeWriteFileError, err.Error(), err)
	}

	client := utils.NewSsh(
		executor.Account.AccountName,
		password,
		executor.IPAddr,
		executor.SSHPort)

	err = client.TransferFile(srcPath, dstPath)
	if err != nil {
		return serializer.Err(serializer.CodeTransferFileError, err.Error(), err)
	}

	var (
		errChan = make(chan error)
		outChan = make(chan string)
	)

	go utils.RunSafeWithMsg(func() {
		out, err := client.RemoteCommand(command)
		if err != nil {
			errChan <- err
		}
		outChan <- out
	}, "调试脚本出错")
	select {
	case <-ctx.Done():
		return serializer.Err(http.StatusGatewayTimeout, "脚本执行超时", nil)
	case err := <-errChan:
		return serializer.Err(serializer.CodeTestShellError, err.Error(), err)
	case out := <-outChan:
		return &serializer.Response{Data: out}
	}
}

func (ss *ScriptService) UpdateScript() *serializer.Response {

	// 判断脚本是否存在危险命令
	err, needAudit := CheckDangerCmd(ss.Script.ScriptContent)
	if err != nil {
		return serializer.Err(serializer.CodeDangerCmdQueryError, "校验危险命名错误", err)
	}

	if needAudit {
		ss.Script.Status = consts.NoAudit
	} else {
		ss.Script.Status = consts.PassAudit
	}

	err = mapper.GetScriptMapper().UpdateScript(ss.Script, ss.UserId)
	if err != nil {
		return serializer.DBErr("更新脚本失败", err)
	}
	return &serializer.Response{Data: ss.Script, Message: "ok"}
}

func (ss *ScriptService) AddScript() *serializer.Response {
	// 查询脚本名是否存在
	scriptMapper := mapper.GetScriptMapper()
	script, err := scriptMapper.FindByName(ss.ScriptName)
	if err != nil {
		return serializer.DBErr("查询脚本失败", err)
	}
	if script != nil {
		return serializer.Err(serializer.CodeScriptAlreadyExist, "脚本名称已存在", err)
	}
	// 判断脚本是否存在危险命令
	err, needAudit := CheckDangerCmd(ss.Script.ScriptContent)
	if err != nil {
		return serializer.Err(serializer.CodeDangerCmdQueryError, "校验危险命名错误", err)
	}

	if !needAudit {
		ss.Script.Status = consts.PassAudit
	} else {
		ss.Script.Status = consts.NoAudit

		ss.Script.ScriptAudit = &models.ScriptAudit{
			Reviewer: ss.Script.LastOperator,
			UserRef:  ss.UserId,
		}
	}

	if err := scriptMapper.AddScript(ss.Script); err != nil {
		return serializer.DBErr("新增脚本失败", err)
	}

	return &serializer.Response{Data: ss.Script, Message: "ok"}
}

// CheckShellScript 校验shell类型脚本
func (ss *ScriptService) CheckShellScript(script string) *serializer.Response {
	tempName := utils.BuilderStr("/tmp/", utils.NewUuid(), ".sh")

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

func (rs *RetrieveScriptService) RetrieveScript() *serializer.Response {
	scriptMapper := mapper.GetScriptMapper()
	filter := &models.Script{
		BaseModel: models.BaseModel{
			ID: rs.ID,
		},
	}
	script, err := scriptMapper.FindOne(filter)
	if err != nil {
		logger.Error("查询脚本出错, err:%s", err.Error())
		return serializer.DBErr("查询脚本出错", err)
	}
	return &serializer.Response{Data: script}
}

func (rs *RetrieveScriptService) DeleteScript() *serializer.Response {
	scriptMapper := mapper.GetScriptMapper()
	filter := &models.Script{
		BaseModel: models.BaseModel{
			ID: rs.ID,
		},
	}
	_, err := scriptMapper.Delete(filter)
	if err != nil {
		logger.Error("删除脚本出错, err:%s", err.Error())
		return serializer.DBErr("删除脚本出错", err)

	}
	return &serializer.Response{Message: "ok"}
}

func (ls *ListService) ScriptsList() (count int, rows interface{}, err error) {

	var (
		scriptMapper = mapper.GetScriptMapper()
		filter       = &models.Script{}
		scripts      = &[]*models.Script{}
	)

	ls.ValidDate()
	count, err = scriptMapper.Count(filter, ls.Sort, ls.Conditions, ls.Searches)
	if err != nil {
		logger.Error("查询标签总数失败: [%s]", err.Error())
		return count, scripts, err
	}
	_, err = scriptMapper.FindAllWithPager(filter, scripts, ls.PageSize, ls.PageNo,
		ls.Sort, ls.Conditions, ls.Searches)

	if err != nil {
		logger.Error("查询标签列表失败: [%s]", err.Error())
		return count, scripts, err
	}
	return count, scripts, err
}

func CheckDangerCmd(scriptContent string) (error, bool) {
	// 判断脚本是否存在危险命令
	cmdFilter := &models.DangerousCmd{}
	commands, err := mapper.GetDangerCmdMapper().ListAllDangerousCommand(cmdFilter)
	if err != nil {
		return err, false
		//return serializer.DBErr("查询危险命令出错", err), false
	}
	var needAudit bool
	for _, c := range commands {
		if c.CheckExistInScript(scriptContent) {
			needAudit = true
		}
	}
	return nil, needAudit
}
