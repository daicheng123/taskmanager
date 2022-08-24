package mapper

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"taskmanager/internal/models"
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

//FindAll 查询所有满足条件的：不包括软删除的
func (bm *BaseMapper) FindAll(filter interface{}, results interface{}) (*gorm.DB, error) {
	return store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Where(filter).Find(results)
	})
}

//Save Create Or Update  创建或更新对象
func (bm *BaseMapper) Save(conflictKeys []clause.Column, value models.UniqKeyGenerator, omitColumns ...string) error {
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		tx := db.Session(&gorm.Session{FullSaveAssociations: true}).Clauses(clause.OnConflict{
			UpdateAll: true,
			Columns:   conflictKeys,
		})
		bm.lock()
		defer bm.Unlock()

		if len(omitColumns) > 0 {
			return tx.Omit(omitColumns...).Create(value)
		}
		return tx.Create(value)
	})
	return err
}
