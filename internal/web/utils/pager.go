package utils

// PagerResult 分页响应
type PagerResult struct {
	PageSize  int         `json:"pageSize"`
	PageNo    int         `json:"pageNo"`
	Count     int         `json:"count"`
	TotalPage int         `json:"totalPage"`
	Rows      interface{} `json:"rows"`
}

func (p *PagerResult) CompletePageInfo() {
	if p.PageSize == 0 {
		return
	}
	mod := p.Count % p.PageSize
	if mod > 0 {
		p.TotalPage = p.Count/p.PageSize + 1
	} else {
		p.TotalPage = p.Count / p.PageSize
	}

	if p.PageNo > p.TotalPage {
		p.PageNo = p.TotalPage
	}
}
