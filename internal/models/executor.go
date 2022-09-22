package models

import (
	"fmt"
	"gorm.io/plugin/soft_delete"
)

type Executor struct {
	*BaseModel
	HostName     string                `json:"hostName"     gorm:"column:hostName;string;size:128;not null;uniqueIndex:idx_hn;commit:执行器主机名"`
	IPAddr       string                `json:"ipAddr"       gorm:"column:ipAddr;string;size:15; not null;commit:执行器ip地址"`
	SHHPort      uint                  `json:"sshPort"      gorm:"column:shhPort;uint;size:5; not null;commit:ssh端口"`
	Status       uint                  `json:"status"       gorm:"column:status;uint;size:1;not null;default:1;commit:执行器状态,1未知 2正常 3不可达"`
	SecretStatus uint                  `json:"secretStatus" gorm:"column:secretStatus;uint;size:1;not null;default:1;commit:密钥分发状态, 1未分发  2分发中 3已分发 4分发失败"`
	Accounts     *ExecutorAccount      `json:",inline"      gorm:"foreignKey:ExecutorRef;references:ID"`
	ExecutePath  string                `json:"executePath"  gorm:"column:executePath;string;size:128;default:/opt/taskmanager;not null;commit:任务执行路径"`
	LastOperator string                `json:"lastOperator" gorm:"type:string; column:lastOperator; size:32; not null;commit:最后操作人"`
	Remarks      string                `json:"remarks"      gorm:"column:remarks;text;not null;commit:备注"`
	DeletedAt    soft_delete.DeletedAt `json:"-" gorm:"index;column:deletedAt;uniqueIndex:idx_hn"`
}

func (em *Executor) TableName() string {
	return "executors"
}

func (em *Executor) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", em.HostName, em.ID)
}

type ExecutorAccount struct {
	BaseModel
	AccountName     string                `json:"accountName" gorm:"column:accountName;string;size:32;not null;uniqueIndex:idx_eid_n;commit:账号名称"`
	ExecutorRef     uint                  `gorm:"column:executorRef;uint;not null;uniqueIndex:idx_eid_n;commit:关联执行器"`
	AccountPassword string                `json:"accountPassword" gorm:"column:accountPassword;string:256;not null;commit:账号密码"`
	DeletedAt       soft_delete.DeletedAt `gorm:"index;column:deletedAt;uniqueIndex:idx_eid_n" json:"-"`
}

func (eam *ExecutorAccount) TableName() string {
	return "executor_accounts"
}

func (eam *ExecutorAccount) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", eam.AccountName, eam.ID)
}
