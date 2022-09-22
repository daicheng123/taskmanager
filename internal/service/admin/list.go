package admin

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
