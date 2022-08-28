package models

import (
	"fmt"
)

type UserModel struct {
	BaseModel
	UserName     string `gorm:"not null; column:userName; type:varchar(128)" json:"userName"`
	UserCode     string `gorm:"not null; column:userCode; type:varchar(128);uniqueIndex:idx_un_code" json:"userCode"`
	UserPassword string `gorm:"not null; column:userPassword; type:varchar(128)" json:"-"`
	UserEmail    string `gorm:"not null; column:userMail; type:varchar(128);uniqueIndex:idx_un_email" json:"email"`
}

func (um *UserModel) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%s", um.ID, um.UserCode)
}

func (um *UserModel) TableName() string {
	return "users"
}
