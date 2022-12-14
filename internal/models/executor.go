package models

import (
	"fmt"
	"gorm.io/plugin/soft_delete"
	"taskmanager/internal/consts"
)

type Executor struct {
	BaseModel
	HostName     string                `gorm:"column:hostName;string;size:128;not null;uniqueIndex:idx_hn;commit:执行器主机名"`
	IPAddr       string                `gorm:"column:ipAddr;string;size:15; not null;uniqueIndex:idx_ip;commit:执行器ip地址"`
	SSHPort      uint                  `gorm:"column:sshPort;uint;size:5; not null;commit:ssh端口"`
	Status       consts.ExecutorStatus `gorm:"column:status;uint;size:1;not null;default:1;commit:执行器状态,1未知 2正常 3不可达"`
	SecretStatus uint                  `gorm:"column:secretStatus;uint;size:1;not null;default:1;commit:密钥分发状态, 1未分发  2分发中 3已分发 4分发失败"`
	Account      *ExecutorAccount      `gorm:"foreignKey:ExecutorRef;references:ID"`
	ExecutePath  string                `gorm:"column:executePath;string;size:128;default:/opt/taskmanager;not null;commit:任务执行路径"`
	LastOperator string                `gorm:"type:string; column:lastOperator; size:32; not null;commit:最后操作人"`
	Remarks      string                `gorm:"column:remarks;text;not null;commit:备注"`
	DeletedAt    soft_delete.DeletedAt `gorm:"index;column:deletedAt;uniqueIndex:idx_hn,idx_ip"`
}

func (em *Executor) TableName() string {
	return "executors"
}

func (em *Executor) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", em.HostName, em.ID)
}

type ExecutorAccount struct {
	BaseModel
	AccountName     string                `gorm:"column:accountName;string;size:32;not null;uniqueIndex:idx_eid_n;commit:账号名称"`
	ExecutorRef     uint                  `gorm:"column:executorRef;uint;not null;uniqueIndex:idx_eid_n;commit:关联执行器"`
	AccountPassword string                `gorm:"column:accountPassword;string:256;not null;commit:账号密码"`
	DeletedAt       soft_delete.DeletedAt `gorm:"index;column:deletedAt;uniqueIndex:idx_eid_n"`
}

func (eam *ExecutorAccount) TableName() string {
	return "executor_accounts"
}

func (eam *ExecutorAccount) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", eam.AccountName, eam.ID)
}
