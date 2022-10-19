package models

import (
	"fmt"
)

type Tag struct {
	BaseModel
	TagName      string `gorm:"type:string; column:tagName; size:256; not null; uniqueIndex:idx_un_tn" json:"tagName"`
	LastOperator string `gorm:"type:string; column:lastOperator; size:32; not null;" json:"lastOperator"`
	// 级联约束条件下不能使用软删除 https://so.muouseo.com/qa/d109nnedx5e4.html
	//DeletedAt    soft_delete.DeletedAt `gorm:"index; column:deletedAt; uniqueIndex:idx_un_tn" json:"-"`
}

func (tm *Tag) GenerateUniqKey() string {
	return fmt.Sprintf("%s_%d", tm.TagName, tm.ID)
}

func (tm *Tag) TableName() string {
	return "tags"
}
