package admin

import (
	"taskmanager/internal/mapper"
	"taskmanager/internal/models"
	"taskmanager/internal/web/utils"
	"taskmanager/pkg/logger"
	"taskmanager/utils/serializer"
)

type TagsService struct {
	TagName      string `json:"tagName"      binding:"required"`
	LastOperator string `json:"lastOperator" binding:"required"`
	*ListService `form:", inline"`
}

func (ts *TagsService) TagSave() serializer.Response {
	tag := &models.TagsModel{
		TagName:      ts.TagName,
		LastOperator: ts.LastOperator,
	}

	if err := mapper.GetTagsMapper().Save(tag); err != nil {
		logger.Error("新增标签失败: [%s]", err.Error())
		return serializer.DBErr("新增标签失败", err)
	}
	return serializer.Response{}
}

func (ts *TagsService) ListByPager() serializer.Response {
	ts.ValidDate()
	filter := &models.TagsModel{}
	tags := &[]*models.TagsModel{}

	count, err := mapper.GetTagsMapper().Count(filter, ts.SortBy, ts.Conditions, ts.Searches)
	if err != nil {
		logger.Error("查询标签总数失败: [%s]", err.Error())
		return serializer.DBErr("查询标签列表失败", err)
	}
	_, err = mapper.GetTagsMapper().FindAllWithPager(filter, tags, ts.PageSize, ts.PageNo,
		ts.SortBy, ts.Conditions, ts.Searches)

	if err != nil {
		logger.Error("查询标签列表失败: [%s]", err.Error())
		return serializer.DBErr("查询标签列表失败", err)
	}

	result := &utils.PagerResult{
		PageSize: ts.PageSize,
		PageNo:   ts.PageNo,
		Count:    count,
	}

	result.CompletePageInfo()
	result.Rows = tags
	return serializer.Response{Data: result}
}
