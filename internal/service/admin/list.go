package admin

type ListService struct {
	PageSize   int    `form:"pageSize" binding:"omitempty"`
	PageNo     int    `form:"pageNo"   binding:"omitempty"`
	OrderBy    string `form:"orderBy"  binding:"omitempty"`
	Order      string `form:"order"    binding:"omitempty,required,eq=DESC|eq=ASC"`
	Sort       string
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

	if len(ls.Order) == 0 {
		ls.Order = "DESC"
	}

	if len(ls.OrderBy) == 0 {
		ls.OrderBy = "updatedAt"
	}
	ls.Sort = ls.OrderBy + " " + ls.Order
}
