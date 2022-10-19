package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"taskmanager/internal/dal/mapper"
	"taskmanager/internal/models"
	"taskmanager/internal/models/common"
	webutils "taskmanager/internal/web/utils"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
)

type AccountService struct {
	AccountName     string `json:"accountName" form:"accountName" binding:"required"`
	AccountPassword string `json:"accountPassword" form:"accountPassword" binding:"required"`
}

type ExecutorService struct {
	ID              uint              `json:"id" uri:"id"  binding:"omitempty,gte=1"`
	HostName        string            `json:"hostName" form:"hostName" binding:"omitempty"`
	IPAddr          string            `json:"ipAddr"   form:"ipAddr" binding:"required"`
	SSHPort         uint              `json:"sshPort"  form:"sshPort" binding:"required,min=1,max=65535"`
	ExecutePath     string            `json:"executePath"`
	ShouldDelivery  bool              `json:"shouldDelivery"`
	SecretStatus    uint              `json:"secretStatus"`
	Status          uint              `json:"status" binding:"omitempty,oneof=0 1 2"`
	LastOperator    string            `json:"lastOperator"`
	UpdatedAt       common.CustomTime `json:"updatedAt"`
	*AccountService `json:",inline"  form:",inline" binding:"required"`
	Remarks         string `json:"remarks"`
}
type ExecutorDelService struct {
	ID uint `uri:"id" binding:"omitempty,gte=1"`
}

type ExeTestService struct {
	ID              uint   `form:"id" binding:"omitempty" `
	IPAddr          string `form:"ipAddr" binding:"required"`
	SSHPort         uint   `form:"sshPort" binding:"required,min=1,max=65535"`
	AccountName     string `form:"accountName" binding:"required"`
	AccountPassword string `form:"accountPassword" binding:"required"`
}

func (ets *ExeTestService) RefreshNode(flag string) *serializer.Response {
	var (
		status uint
		resp   *serializer.Response
	)
	fn := func(status uint) error {
		value := &models.Executor{
			BaseModel: models.BaseModel{
				ID: ets.ID,
			},
			Status: status,
		}
		return mapper.NewExecutorMapper().Updates(value)
	}

	if err := testNode(ets.AccountName, ets.AccountPassword, ets.IPAddr, ets.SSHPort); err != nil {
		//status = 3
		status = utils.HostUnreachable
		resp = serializer.Err(serializer.CodeServerInternalError, err.Error(), err)
	} else {
		//status = 2
		status = utils.HostAvail
		resp = &serializer.Response{}
	}
	if flag == "update" {
		go utils.RunSafeWithMsg(func() {
			err := fn(status)
			if err != nil {
				logger.Error("更新执行器状态失败, err:[%s]", err.Error())
			}
		}, "更新执行器状态失败")
	}
	return resp
}

func (es *ExecutorService) ExecutorOptions() *serializer.Response {
	filter := &models.Executor{
		Status:       utils.HostAvail,
		SecretStatus: utils.KeyDistributed,
	}
	executors := make([]models.Executor, 0)
	_, err := mapper.GetExecutorMapper().FindAll(filter, &executors)
	if err != nil {
		return serializer.DBErr("获取执行器失败", err)
	}

	for _, executor := range executors {
		fmt.Printf("id:%dhahahahahaha\n", executor.ID)
	}
	return &serializer.Response{Data: executors}
}

// ExecutorAdd 创建执行器
func (es *ExecutorService) ExecutorAdd() *serializer.Response {
	em, err := ExecutorTransToModel(es)
	if err != nil {
		return serializer.Err(http.StatusInternalServerError, err.Error(), err)
	}
	err = mapper.GetExecutorMapper().Save(em)
	if err != nil {
		logger.Error("执行器新增失败, err:[%s]", err.Error())
		return serializer.DBErr("执行器新增失败", err)
	}

	if es.ShouldDelivery {
		return es.DistributeKey(em)
	}
	return &serializer.Response{Data: em}
}

// ExecutorUpdate 更新执行器
func (es *ExecutorService) ExecutorUpdate() *serializer.Response {
	filter := &models.Executor{
		BaseModel: models.BaseModel{
			ID: es.ID,
		},
	}
	exeObj, err := mapper.GetExecutorMapper().PreLoadFindOne(filter)
	if err != nil {
		logger.Error("查询执行器失败,ID: %d, err:[%s]", es.ID, err.Error())
		return serializer.DBErr("查询执行器失败", err)
	}
	newPwd, err := utils.Encrypt(es.AccountPassword)
	if err != nil {
		return serializer.Err(http.StatusInternalServerError, "密码解析失败", err)
	}

	exeObj.HostName = es.HostName
	exeObj.ExecutePath = es.ExecutePath
	exeObj.Accounts.AccountPassword = newPwd
	exeObj.Status = es.Status
	exeObj.Remarks = es.Remarks
	exeObj.LastOperator = es.LastOperator
	// 用户名/端口号变化后则主机密钥分发重置
	if exeObj.Accounts.AccountName != es.AccountName || exeObj.SHHPort != es.SSHPort {
		//exeObj.SecretStatus = 1
		exeObj.SecretStatus = utils.KeyUndistributed
	}
	exeObj.Accounts.AccountName = es.AccountName
	exeObj.SHHPort = es.SSHPort

	err = mapper.GetExecutorMapper().Updates(exeObj)
	if err != nil {
		logger.Error("更新执行器失败,ID: %d, err:[%s]", es.ID, err.Error())
		return serializer.DBErr("更新执行器失败", err)
	}
	return &serializer.Response{Data: exeObj, Message: "ok"}
}

// DistributeKey 分发主机密钥
func (es *ExecutorService) DistributeKey(executor *models.Executor) *serializer.Response {
	executorMapper := mapper.GetExecutorMapper()
	//executor.SecretStatus = 2
	executor.SecretStatus = utils.KeyDistributing
	if err := executorMapper.Updates(executor); err != nil {
		return serializer.DBErr("更新执行器密钥状态失败", err)
	}

	sshCli := utils.NewSsh(es.AccountName, es.AccountPassword, es.IPAddr, es.SSHPort)
	if _, err := sshCli.DistributeKey(); err != nil {
		//executor.SecretStatus = 4
		executor.SecretStatus = utils.KeyDistributeFailed
		go utils.RunSafeWithMsg(func() {
			err = executorMapper.Updates(executor)
			if err != nil {
				logger.Error("未能将钥状态更新为分发失败状态, err:[%s]", err.Error())
			}
		}, "未能将钥状态更新为分发失败状态")
		return serializer.Err(serializer.CodeServerInternalError, "", err)
	}

	//executor.SecretStatus = 3
	executor.SecretStatus = utils.KeyDistributed
	if err := executorMapper.Updates(executor); err != nil {
		go utils.RunSafeWithMsg(func() {
			err = executorMapper.Updates(executor)
			if err != nil {
				logger.Error("更新执行器密钥状态失败")
			}
		}, "回滚密钥状态失败")
	}
	return &serializer.Response{Message: "ok"}
}

func (ed *ExecutorDelService) ExecutorDelete() *serializer.Response {
	filter := &models.Executor{
		BaseModel: models.BaseModel{
			ID: ed.ID,
		},
	}
	_, err := mapper.GetExecutorMapper().Delete(filter)
	if err != nil {
		return serializer.DBErr("删除执行器失败", err)
	}
	return &serializer.Response{}
}

func (ed *ExecutorDelService) ExecutorBatchDelete(ids []uint) *serializer.Response {
	// 批量删除
	err := mapper.GetExecutorMapper().BatchDeleteById(ids...)
	if err != nil {
		logger.Error("执行器批量删除失败,err:[%s]", err.Error())
		return serializer.DBErr("执行器批量删除失败", err)
	}
	return &serializer.Response{Data: ids}
}

// ExecutorList list 执行器
func (ls *ListService) ExecutorList() *serializer.Response {
	ls.ValidDate()
	filter := &models.Executor{}
	executors := &[]*models.Executor{}
	count, err := mapper.GetExecutorMapper().Count(filter, ls.Sort, ls.Conditions, ls.Searches)
	if err != nil {
		logger.Error("查询执行器总数失败: [%s]", err.Error())
		return serializer.DBErr("查询执行器总数失败", err)
	}
	err = mapper.GetExecutorMapper().FindAllWithPager(filter, executors, ls.PageSize, ls.PageNo,
		ls.Sort, ls.Conditions, ls.Searches)
	if err != nil {
		logger.Error("查询执行器列表失败: [%s]", err.Error())
		return serializer.DBErr("查询执行器列表失败", err)
	}
	result := webutils.PagerResult{
		PageSize: ls.PageSize,
		PageNo:   ls.PageNo,
		Count:    count,
	}
	result.CompletePageInfo()
	result.Rows, err = ExecutorTransServices(*executors)
	if err != nil {
		return serializer.Err(http.StatusInternalServerError, "解析executors列表失败", err)
	}
	return &serializer.Response{Data: result}
}

func testNode(name, password, ip string, port uint) error {
	sshCli := utils.NewSsh(name, password, ip, port)
	stdOut, err := sshCli.RemoteCommand("echo success")
	if err != nil {
		return err
	}
	if stdOut != "" || strings.Contains(stdOut, "success") {
		return nil
	}
	return errors.New("执行测试失败")
}

// ExecutorTransToModel 执行器dto
func ExecutorTransToModel(executor *ExecutorService) (*models.Executor, error) {
	accountPwd, err := utils.Encrypt(executor.AccountPassword)
	if err != nil {
		logger.Error("主机密码解析错误, err:[%s]", err.Error())
		return nil, serializer.NewError(serializer.CodeHostPasswordEncodeErr, "主机密码解析错误", err)
	}

	accounts := &models.ExecutorAccount{
		AccountName:     executor.AccountName,
		AccountPassword: accountPwd,
	}
	return &models.Executor{
		HostName:     executor.HostName,
		IPAddr:       executor.IPAddr,
		SHHPort:      executor.SSHPort,
		LastOperator: executor.LastOperator,
		Status:       executor.Status,
		ExecutePath:  executor.ExecutePath,
		Accounts:     accounts,
	}, nil
}

func ExecutorTransServices(executors []*models.Executor) (services []*ExecutorService, err error) {
	services = make([]*ExecutorService, 0)
	for _, e := range executors {
		password, err := utils.Decrypt(e.Accounts.AccountPassword)
		if err != nil {
			logger.Error("执行器:%s 密码解析错误,err:[%s]", e.HostName, err.Error())
			break
		}
		es := &ExecutorService{
			ID:           e.ID,
			HostName:     e.HostName,
			IPAddr:       e.IPAddr,
			SSHPort:      e.SHHPort,
			ExecutePath:  e.ExecutePath,
			LastOperator: e.LastOperator,
			Status:       e.Status,
			SecretStatus: e.SecretStatus,
			UpdatedAt:    e.UpdatedAt,
			Remarks:      e.Remarks,
			AccountService: &AccountService{
				AccountName:     e.Accounts.AccountName,
				AccountPassword: password,
			},
		}
		services = append(services, es)
	}
	return services, err
}
