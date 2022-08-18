package mapper

import (
	"gorm.io/gorm"
	"taskmanager/pkg/store"
)

type BaseMapper struct {
	Lock   func()
	Unlock func()
	locked bool
}

func (bm *BaseMapper) lock() {
	if bm.Lock != nil {
		bm.Lock()
		bm.locked = true
	}
}

func (bm *BaseMapper) unlock() {
	if bm.locked && bm.Unlock != nil {
		bm.Unlock()
	}
}

// FindOne 获取满足条件的第一条记录(不排序)
func (bm *BaseMapper) FindOne(filter, result interface{}) (*gorm.DB, error) {
	return store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Where(filter).Take(result)
	})
}
