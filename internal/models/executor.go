package models

import (
	"fmt"
	"gorm.io/plugin/soft_delete"
)

type Executor struct {
	BaseModel
	HostName     string                `gorm:"column:hostName;string;size:128;not null;uniqueIndex:idx_hn;commit:执行器主机名"`
	IPAddr       string                `gorm:"column:ipAddr;string;size:15; not null;commit:执行器ip地址"`
	SHHPort      uint                  `gorm:"column:shhPort;uint;size:5; not null;commit:ssh端口"`
	Status       uint                  `gorm:"column:status;uint;size:1;not null;commit:执行器状态,0未知 1正常 2不可达"`
	SecretStatus uint                  `gorm:"column:secretStatus;uint;size:1;not null;default:0;commit:密钥分发状态, 0未分发 1已分发 3分发中 4不包含root账号 5密码错误 6主机类型错误"`
	Account      *ExecutorAccount      `gorm:"foreignKey:ExecutorRef;references:ID"`
	ExecutePath  string                `gorm:"column:executePath;string;size:128;not null;commit:任务执行路径"`
	LastOperator string                `gorm:"type:string; column:lastOperator; size:32; not null;commit:最后操作人"`
	Remarks      string                `gorm:"column:remarks;string;size:512;not null;commit:备注"`
	DeletedAt    soft_delete.DeletedAt `gorm:"index;column:deletedAt;uniqueIndex:idx_hn" json:"-"`
}

func (em *Executor) TableName() string {
	return "executors"
}

func (em *Executor) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%s", em.ID, em.HostName)
}

type ExecutorAccount struct {
	BaseModel
	AccountName     string `gorm:"column:accountName;string;size:32;not null;uniqueIndex:idx_eid_n;commit:账号名称"`
	ExecutorRef     uint   `gorm:"column:executorRef;uint;not null;uniqueIndex:idx_eid_n;commit:关联执行器"`
	AccountPassword string `gorm:"column:accountPassword;string:256;not null;commit:账号密码"`
	//ConnectWay      uint                  `gorm:"column:connectWay;uint;default:1;not null;commit:连接方式, 1ssh"`
	//Remarks         string                `gorm:"column:remarks;string;commit:备注"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;column:deletedAt;uniqueIndex:idx_eid_n" json:"-"`
}

func (eam *ExecutorAccount) TableName() string {
	return "executor_accounts"
}

func (eam *ExecutorAccount) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%s", eam.ID, eam.AccountName)
}
