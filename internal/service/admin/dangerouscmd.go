package admin

import (
	"taskmanager/internal/dal/mapper"
	"taskmanager/internal/models"
	"taskmanager/internal/web/utils"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
)

type DangerousCommandService struct {
	ID      uint   `json:"id" binding:"omitempty,gte=1"`
	Command string `json:"dangerousCommand" binding:"required,min=2,max=100"`
	Remarks string `json:"remarks"`
}

type DangerousCommandDelService struct {
	ID uint `uri:"id" binding:"omitempty,gte=1"`
}

func (dcs *DangerousCommandService) DangerousCmdSave() *serializer.Response {
	dangerCmd := &models.DangerousCmd{
		Command: dcs.Command,
		Remarks: dcs.Remarks,
	}
	// 更新
	if dcs.ID != 0 {
		dangerCmd.BaseModel = &models.BaseModel{
			ID: dcs.ID,
		}
	}
	err := mapper.GetDangerCmdMapper().Upsert(dangerCmd)
	if err != nil {
		return serializer.DBErr("危险命令保存失败", err)
	}
	return &serializer.Response{Message: "ok", Data: dangerCmd}
}

func (dcd *DangerousCommandDelService) DangerousCmdDelete() *serializer.Response {
	filter := &models.DangerousCmd{
		BaseModel: &models.BaseModel{
			ID: dcd.ID,
		},
	}
	_, err := mapper.GetDangerCmdMapper().Delete(filter)
	if err != nil {
		logger.Error("删除危险命令失败, err:[%s]", err.Error())
		return serializer.DBErr("删除危险命令失败", err)
	}
	return &serializer.Response{Message: "ok"}
}

func (ls *ListService) DangerousCmdList() *serializer.Response {
	ls.ValidDate()
	filter := &models.DangerousCmd{}
	commands := &[]models.DangerousCmd{}

	count, err := mapper.GetDangerCmdMapper().Count(filter, ls.Sort)
	if err != nil {
		logger.Error("查询危险命令总数失败: [%s]", err.Error())
		return serializer.DBErr("查询危险命令表失败", err)
	}
	err = mapper.GetDangerCmdMapper().FindAllWithPager(filter, commands, ls.PageSize, ls.PageNo)

	if err != nil {
		logger.Error("查询危险命令列表失败: [%s]", err.Error())
		return serializer.DBErr("查询危险命令列表失败", err)
	}
	result := &utils.PagerResult{
		PageSize: ls.PageSize,
		PageNo:   ls.PageNo,
		Count:    count,
	}
	result.CompletePageInfo()
	result.Rows = commands
	return &serializer.Response{Data: result}
}
