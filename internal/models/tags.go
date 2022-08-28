package models

import (
	"fmt"
	"gorm.io/plugin/soft_delete"
)

type TagsModel struct {
	BaseModel
	TagName      string                `gorm:"type:string; column:tagName; size:256; not null; uniqueIndex:idx_un_tn" json:"tagName"`
	LastOperator string                `gorm:"type:string; column:lastOperator; size:32; not null;" json:"lastOperator"`
	DeletedAt    soft_delete.DeletedAt `gorm:"index; column:deletedAt; uniqueIndex:idx_un_tn" json:"-"`
}

func (tm *TagsModel) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%s", tm.ID, tm.TagName)
}

func (tm *TagsModel) TableName() string {
	return "tags"
}
