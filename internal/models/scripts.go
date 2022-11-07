package models

import (
	"fmt"
	"taskmanager/internal/consts"
)

type Script struct {
	BaseModel
	ScriptName    string              `json:"scriptName" gorm:"column:scriptName; string;size:128;not null;uniqueIndex:idx_sn;commit:脚本名称"`
	ScriptContent string              `json:"scriptContent" gorm:"column:scriptContent; text;not null;size:500;commit:脚本内容"`
	Status        consts.ScriptStatus `json:"status" gorm:"column:status;uint;size:1;not null;default:1;commit:脚本审核状态,1未审核 2审核中 3已通过 4驳回"`
	ScriptType    consts.ScriptType   `json:"scriptType" gorm:"column:scriptType;uint;size:1;not null;commit:脚本类型,1shell 2python"`
	OverTime      uint                `json:"overTime" gorm:"column:overTime;uint;size:4;not null;commit:脚本超时时间"`
	Tag           *Tag                `json:"tag,omitempty" gorm:"foreignKey:TagRefer; references:ID; constraint:OnDelete:RESTRICT"`
	TagRefer      uint                `json:"tagId" gorm:"column:tagRefer"`
	ScriptAudit   *ScriptAudit        `gorm:"foreignKey:ScriptRef; references:ID; constraint:OnDelete:CASCADE"`
	LastOperator  string              `json:"lastOperator" gorm:"column:lastOperator;type:string;size:32; not null;commit:最后操作人"`
	Remarks       string              `json:"remarks" gorm:"column:remarks;text;not null;commit:备注"`
	//DeletedAt     soft_delete.DeletedAt `gorm:"uniqueIndex:idx_sn"`
}

func (s *Script) TableName() string {
	return "scripts"
}

func (s *Script) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", s.ScriptName, s.ID)
}

type ScriptAudit struct {
	BaseModel
	ScriptRef uint       `json:"scriptId" gorm:"column:scriptRef"`
	UserRef   uint       `json:"userRef" gorm:"column:userRef; commit:提交人"`
	Applicant *UserModel `gorm:"foreignKey:UserRef; references:ID; constraint:OnDelete:SET NULL"`
	Reviewer  string     `json:"reviewer" gorm:"column:reviewer;type:string;size:32; not null;commit:审核人"`
}

func (s *ScriptAudit) TableName() string {
	return "script_audit"
}

func (s *ScriptAudit) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%d", s.ScriptRef, s.ID)
}
