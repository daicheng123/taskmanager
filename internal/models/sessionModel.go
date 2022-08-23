package models

import (
	"fmt"
	"time"
)

type SessionModel struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"type:uint; size:12; not null; uniqueIndex:idx_un_ui" json:"userId"`
	SessionValue string    `gorm:"type:string; size:256; not null; uniqueIndex:idx_un_sv"`
	ExpireTime   time.Time `gorm:"type:time; not null"`
}

func (sm *SessionModel) GenerateUniqKey() string {
	return fmt.Sprintf("%d_%s", sm.ID, sm.SessionValue)
}

func (sm *SessionModel) TableName() string {
	return "login_session"
}
