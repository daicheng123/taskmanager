package transfers

import (
	"taskmanager/internal/models"
	"taskmanager/internal/schemas"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
)

func ExecutorModel2Data(executor *models.Executor) (*schemas.ExecutorInfo, error) {
	accountPwd, err := utils.Decrypt(executor.Account.AccountPassword)
	if err != nil {
		logger.Error("主机密码解析错误, err:[%s]", err.Error())
		return nil, serializer.NewError(serializer.CodeHostPasswordDecodeErr, "主机密码解析错误", err)
	}

	var executorInfo = &schemas.ExecutorInfo{
		ID:           executor.ID,
		CreatedAt:    executor.CreatedAt,
		UpdatedAt:    executor.UpdatedAt,
		HostName:     executor.HostName,
		IPAddr:       executor.IPAddr,
		SSHPort:      executor.SSHPort,
		Status:       executor.Status,
		SecretStatus: executor.SecretStatus,
		AccountInfo: &schemas.AccountInfo{
			AccountName:     executor.Account.AccountName,
			AccountPassword: accountPwd,
		},
		ExecutePath:  executor.ExecutePath,
		LastOperator: executor.LastOperator,
	}
	return executorInfo, nil
}

func ExecutorData2Model(executor *schemas.ExecutorInfo) (*models.Executor, error) {
	accountPwd, err := utils.Encrypt(executor.AccountInfo.AccountPassword)
	if err != nil {
		logger.Error("主机密码加密错误, err:[%s]", err.Error())
		return nil, serializer.NewError(serializer.CodeHostPasswordEncodeErr, "主机密码加密错误", err)
	}

	accounts := &models.ExecutorAccount{
		AccountName:     executor.AccountName,
		AccountPassword: accountPwd,
	}

	em := &models.Executor{
		BaseModel: models.BaseModel{
			ID: executor.ID,
		},
		HostName:     executor.HostName,
		IPAddr:       executor.IPAddr,
		SSHPort:      executor.SSHPort,
		Status:       executor.Status,
		SecretStatus: executor.SecretStatus,
		Account:      accounts,
		ExecutePath:  executor.ExecutePath,
		LastOperator: executor.LastOperator,
		Remarks:      executor.Remarks,
	}
	return em, nil
}
