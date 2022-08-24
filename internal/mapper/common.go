package mapper

import (
	"gorm.io/gorm"
	"strings"
)

// 处理分页条件
func paginate(pageSize int, pageNo int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo <= 0 {
			pageNo = 1
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func conditionBy(conditions map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range conditions {
			db = db.Where(k+" = ?", v)
		}
		return db
	}
}

func searchBy(searches map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		search := ""
		for k, v := range searches {
			search += k + " like '%" + v + "%' OR "
		}
		search = strings.TrimSuffix(search, " OR ")
		return db.Where(search)
	}
}

// 处理排序条件
func orderBy(sortBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(sortBy) <= 0 {
			return db
		}
		return db.Order(sortBy)
	}
}
