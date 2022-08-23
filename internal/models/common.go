package models

import (
	"gorm.io/gorm"
	"time"
)

type UniqKeyGenerator interface {
	GenerateUniqKey() string
}

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updateAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
