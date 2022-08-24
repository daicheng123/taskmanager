package models

import (
	"fmt"
)

type TagsModel struct {
	BaseModel
	TagName      string `gorm:"type:string; size:256; not null; uniqueIndex:idx_un_tn" json:"tagName"`
	LastOperator string `gorm:"type:string; size:32; not null;" json:"lastOperator"`
}

func (tm *TagsModel) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%s", tm.ID, tm.TagName)
}

func (tm *TagsModel) TableName() string {
	return "tags"
}
