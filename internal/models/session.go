package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SessionModel struct {
	gorm.Model
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"type:uint; column:userId; size:12; not null; uniqueIndex:idx_un_ui" json:"userId"`
	SessionValue string         `gorm:"type:string; column:sessionValue; size:256; not null; uniqueIndex:idx_un_sv"`
	ExpireTime   time.Time      `gorm:"type:time; column:expireTime; not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index; column:deletedAt" json:"-"`
}

func (sm *SessionModel) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%s", sm.ID, sm.SessionValue)
}

func (sm *SessionModel) TableName() string {
	return "login_session"
}
