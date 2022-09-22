package models

import (
	"fmt"
	"gorm.io/plugin/soft_delete"
	"strings"
)

type DangerousCmd struct {
	*BaseModel
	Command   string                `gorm:"column:command;string;size:100;not null;uniqueIndex:idx_cmd;commit:危险命令" json:"dangerousCmd"`
	Remarks   string                `gorm:"column:remarks;text;not null;commit:备注" json:"remarks"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;column:deletedAt;uniqueIndex:idx_cmd" json:"-"`
}

func (dc *DangerousCmd) TableName() string {
	return "danger_command"
}

func (dc *DangerousCmd) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", dc.Command, dc.ID)
}

func (dc *DangerousCmd) CheckExistInScript(script string) bool {
	return strings.Contains(script, dc.Command)
}
