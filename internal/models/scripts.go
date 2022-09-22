package models

import (
	"fmt"
	"gorm.io/plugin/soft_delete"
)

type Script struct {
	*BaseModel
	ScriptName    string `gorm:"column:scriptName;string;size:128;not null;uniqueIndex:idx_sn;commit:脚本名称"`
	ScriptContent string `gorm:"column:scriptContent;text;not null;size:500;uniqueIndex:idx_sn;commit:脚本内容"`
	Status        uint   `gorm:"column:status;uint;size:1;not null;default:1;commit:脚本审核状态,1未审核 2审核中 3已通过 4驳回"`
	ScriptType    uint   `gorm:"column:status;uint;size:1;not null;commit:脚本类型,1shell 2python"`
	OverTime      uint   `gorm:"column:overTime;uint;size:4;not null;commit:脚本超时时间"`
	Tag           *Tag   `gorm:"foreignKey:TagRefer; references:ID"`
	TagRefer      uint
	LastOperator  string                `gorm:"column:lastOperator;type:string;size:32; not null;commit:最后操作人"`
	Remarks       string                `gorm:"column:remarks;text;not null;commit:备注"`
	DeletedAt     soft_delete.DeletedAt `gorm:"column:deletedAt;uniqueIndex:idx_sn"`
}

func (s *Script) TableName() string {
	return "scripts"
}

func (s *Script) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", s.ScriptName, s.ID)
}
