package models

import (
	"gorm.io/gorm"
	"taskmanager/internal/models/common"
)

type BaseModel struct {
	ID        uint              `gorm:"primarykey" json:"id"`
	CreatedAt common.CustomTime `json:"createdAt"`
	UpdatedAt common.CustomTime `json:"updateAt"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`
}

type UniqKeyGenerator interface {
	GenerateUniqKey() string
}
