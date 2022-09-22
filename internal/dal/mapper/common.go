package mapper

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"taskmanager/utils"
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

// conditionBy 精准匹配
func conditionBy(conditions map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range conditions {
			if !utils.IsZero(v) {
				db = db.Where(k+" = ?", v)
			}
		}
		return db
	}
}

// searchBy 模糊匹配
func searchBy(searches map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		search := ""
		for k, v := range searches {
			if !utils.IsZero(v) {
				search += fmt.Sprintf("%s like '%%%v%%' AND ", k, v)
			}
			//search += k + " like '%" + v + "%' AND "
			//k + " like '%" + v + "%' AND "
		}
		search = strings.TrimSuffix(search, " AND ")
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
