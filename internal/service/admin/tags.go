package admin

import (
	"taskmanager/internal/dal/mapper"
	"taskmanager/internal/models"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	//validator "github.com/go-playground/validator/v10"
)

type TagsService struct {
	ID           uint   `json:"id" binding:"omitempty,gte=1"`
	TagName      string `json:"tagName"      binding:"required"`
	LastOperator string `json:"lastOperator" binding:"required"`
}

func (ts *TagsService) TagSave() *serializer.Response {
	tag := &models.Tag{
		TagName:      ts.TagName,
		LastOperator: ts.LastOperator,
	}

	if err := mapper.GetTagsMapper().Upsert(tag); err != nil {
		logger.Error("新增标签失败: [%s]", err.Error())
		return serializer.DBErr("新增标签失败", err)
	}
	return &serializer.Response{}
}

func (ts *TagsService) TagDelete(id uint) *serializer.Response {
	filter := &models.Tag{
		BaseModel: models.BaseModel{ID: id},
	}
	_, err := mapper.GetTagsMapper().Delete(filter)
	if err != nil {
		logger.Error("删除分类标签失败: [%s]", err.Error())
		return serializer.DBErr("删除分类标签失败", err)
	}
	return &serializer.Response{}
}

func (ts *TagsService) TagEdit() *serializer.Response {
	tag := &models.Tag{
		BaseModel: models.BaseModel{
			ID: ts.ID,
		},
		TagName:      ts.TagName,
		LastOperator: ts.LastOperator,
	}
	if err := mapper.GetTagsMapper().Upsert(tag); err != nil {
		logger.Error("更新标签出错, err:[%s]", err.Error())
		return serializer.DBErr("更新标签出错", err)
	}
	return &serializer.Response{Data: tag}
}

func (ts *TagsService) GetOptions() *serializer.Response {
	filter := &models.Tag{}
	tags, err := mapper.GetTagsMapper().ListAllTags(filter)
	if err != nil {
		logger.Error("查询分类标签失败, err:[%s]", err.Error())
		return serializer.DBErr("查询分类标签失败", err)
	}
	return &serializer.Response{Data: tags}
}

func (ls *ListService) TagsList() (count int, rows interface{}, err error) {
	var (
		tagMapper = mapper.GetTagsMapper()
		filter    = &models.Tag{}
		tags      = &[]*models.Tag{}
	)

	count, err = tagMapper.Count(filter, ls.Sort, ls.Conditions, ls.Searches)
	if err != nil {
		logger.Error("查询标签总数失败: [%s]", err.Error())
		return count, tags, err
	}

	_, err = tagMapper.FindAllWithPager(filter, tags, ls.PageSize, ls.PageNo, ls.Sort, ls.Conditions, ls.Searches)
	if err != nil {
		logger.Error("查询标签列表失败: [%s]", err.Error())
		return count, tags, err
	}

	return count, tags, err
}

//func (ls *ListService) TagsList() *serializer.Response {
//	ls.ValidDate()
//	filter := &models.Tag{}
//	tags := &[]*models.Tag{}
//	testFunc := func() (count int, rows []models.UniqKeyGenerator, err error) {
//		count, err = mapper.GetTagsMapper().Count(filter, ls.Sort, ls.Conditions, ls.Searches)
//		if err != nil {
//			logger.Error("查询标签总数失败: [%s]", err.Error())
//			return count, rows,  err
//		}
//		_, err = mapper.GetTagsMapper().FindAllWithPager(filter, tags, ls.PageSize, ls.PageNo,
//			ls.Sort, ls.Conditions, ls.Searches)
//
//		if err != nil {
//			logger.Error("查询标签列表失败: [%s]", err.Error())
//			return count, rows,  err
//		}
//		return
//	}
//
//	count, rows, err := testFunc()
//	if err != nil {
//		return serializer.DBErr("获取标签数据失败", err)
//	}
//	result := &utils.PagerResult{
//		PageSize: ls.PageSize,
//		PageNo:   ls.PageNo,
//		Count:    count,
//	}
//	result.CompletePageInfo()
//	result.Rows = rows
//	return &serializer.Response{Data: result}
//}
