package mapper

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"taskmanager/internal/models"
	"taskmanager/pkg/store"
	"taskmanager/utils"
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

// PreLoadFindOne
func (bm *BaseMapper) PreLoadFindOne(filter, result interface{}) (*gorm.DB, error) {
	return store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Preload(clause.Associations).Where(filter).Take(result)
	})
}

//FindAll 查询所有满足条件的：不包括软删除的
func (bm *BaseMapper) FindAll(filter interface{}, results interface{}) (*gorm.DB, error) {
	return store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Where(filter).Find(results)
	})
}

//Upsert  创建或更新对象
func (bm *BaseMapper) Upsert(conflictKeys []clause.Column, value models.UniqKeyGenerator, omitColumns ...string) error {
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

func (bm *BaseMapper) Save(value models.UniqKeyGenerator) error {
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		tx := db.Session(&gorm.Session{FullSaveAssociations: true})
		bm.lock()
		defer bm.Unlock()
		return tx.Create(value)
	})
	return err
}

//SoftDeleteByFilter  级联删除
func (bm *BaseMapper) SoftDeleteByFilter(filter interface{}, deletedItems interface{}) (err error) {
	_, err = store.Execute(func(db *gorm.DB) *gorm.DB {
		// 级联删除只能按ID关联，所以这里需要先查询完整的实体
		queryDB, _ := bm.FindAll(filter, deletedItems)
		if utils.Size(deletedItems) <= 0 {
			return queryDB
		}
		bm.lock()
		defer bm.unlock()
		return db.Session(&gorm.Session{}).Select(clause.Associations).Delete(deletedItems)
	})
	return err
}

func (bm *BaseMapper) Updates(value models.UniqKeyGenerator, omitColumns ...string) (int, error) {
	affRows := 0
	if value == nil {
		return 0, nil
	}
	if len(omitColumns) == 0 {
		omitColumns = []string{"deletedAt"}
	}
	db, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		bm.lock()
		defer bm.unlock()
		return db.Session(&gorm.Session{FullSaveAssociations: true}).Model(value).
			Omit(omitColumns...).Updates(value)
	})
	if db != nil {
		affRows = int(db.RowsAffected)
	}
	return affRows, err
}

func (bm *BaseMapper) FindAllWithPager(filter, result interface{}, pageSize, pageNo int,
	sortBy string, conditions, searches map[string]interface{}) (*gorm.DB, error) {
	return store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).
			Model(filter).
			Preload(clause.Associations).
			Scopes(
				conditionBy(conditions),
				searchBy(searches),
				orderBy(sortBy),
				paginate(pageSize, pageNo)).
			Find(result)
	})
}

func (bm *BaseMapper) Count(filter interface{}, sortBy string, conditions, searches map[string]interface{}) (int, error) {
	var count int64
	_, err := store.Execute(func(db *gorm.DB) *gorm.DB {
		return db.Session(&gorm.Session{}).Model(filter).
			Scopes(
				orderBy(sortBy),
				conditionBy(conditions),
				searchBy(searches)).
			Count(&count)
	})
	return int(count), err
}
