package admin

import (
	"taskmanager/internal/web/utils"
	"taskmanager/pkg/serializer"
)

type ListerFunc func() (count int, rows interface{}, err error)

type ListService struct {
	PageSize   int    `form:"pageSize"`
	PageNo     int    `form:"pageNo"`
	OrderBy    string `form:"orderBy"`
	Order      string `form:"order" binding:"omitempty,eq=DESC|eq=ASC"`
	Sort       string
	Conditions map[string]interface{} `form:"conditions"`
	Searches   map[string]interface{} `form:"searches"`
}

func (ls *ListService) ValidDate() {
	if ls.PageSize == 0 {
		ls.PageSize = 10
	}

	if ls.PageNo == 0 {
		ls.PageNo = 1
	}

	if len(ls.Order) == 0 {
		ls.Order = "DESC"
	}

	if len(ls.OrderBy) == 0 {
		ls.OrderBy = "updatedAt"
	}
	ls.Sort = ls.OrderBy + " " + ls.Order
}

func (ls *ListService) Lister(listFunc ListerFunc) *serializer.Response {
	ls.ValidDate()
	count, rows, err := listFunc()
	if err != nil {
		return serializer.DBErr("获取标签数据失败", err)
	}
	result := &utils.PagerResult{
		PageSize: ls.PageSize,
		PageNo:   ls.PageNo,
		Count:    count,
	}
	result.CompletePageInfo()
	result.Rows = rows
	return &serializer.Response{Data: result}
}
