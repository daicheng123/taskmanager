package models

import (
	"taskmanager/internal/models/common"
)

type BaseModel struct {
	ID        uint              `gorm:"primarykey" json:"id"`
	CreatedAt common.CustomTime `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt common.CustomTime `gorm:"column:updatedAt" json:"updatedAt"`
}

type UniqKeyGenerator interface {
	GenerateUniqKey() string
}
