package admin

import "taskmanager/utils/serializer"

type Lister interface {
	ValidDate()
	ListByPager() serializer.Response
}

type ListService struct {
	PageSize   int               `form:"pageSize" binding:"omitempty"`
	PageNo     int               `form:"pageNo"   binding:"omitempty"`
	SortBy     string            `form:"sortBy"   binding:"omitempty"`
	Order      string            `form:"order"    binding:"omitempty,required,eq=DESC|eq=ASC"`
	Conditions map[string]string `form:"conditions" `
	Searches   map[string]string `form:"searches"`
}

func (ls *ListService) ValidDate() {
	if ls.PageSize == 0 {
		ls.PageSize = 10
	}

	if ls.PageNo == 0 {
		ls.PageNo = 1
	}

	if len(ls.SortBy) == 0 {
		ls.SortBy = "created_at"
	}

	if len(ls.SortBy) == 0 {
		ls.SortBy = "DESC"
	}
	//
	//if ls.Conditions != nil && len(ls.Conditions) != 0{
	//
	//}
	//
}

// SortPolicy 默认以 createdAt 进行排序
func (ls *ListService) SortPolicy() string {
	if len(ls.SortBy) > 0 {
		return ls.SortBy
	}
	return "created_at"
}
