package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskmanager/internal/consts"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
	"taskmanager/pkg/logger"
	"taskmanager/utils/serializer"
)

/*
	TagsController 脚本标签管理
*/

type TagsController struct {
}

func NewTagsController() *TagsController {
	return &TagsController{}
}

func (tc *TagsController) tagAdd(ctx *gin.Context) {
	srv := &admin.TagsService{}
	err := ctx.ShouldBindJSON(srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.TagSave())
}

func (tc *TagsController) tagList(ctx *gin.Context) {
	service := &admin.ListService{}
	err := ctx.ShouldBindQuery(service)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, service.TagsList())
}

func (tc *TagsController) tagDel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusOK, serializer.ParamErr("错误的请求参数", nil))
		return
	}
	service := &admin.TagsService{}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("id 解析错误, err:[%s]", err.Error())
		ctx.JSON(http.StatusOK, serializer.Err(serializer.CodeServerInternalError, "服务器内部错误", err))
		return
	}
	ctx.JSON(http.StatusOK, service.TagDelete(uint(id)))
}

func (tc *TagsController) tagEdit(ctx *gin.Context) {
	srv := &admin.TagsService{}
	err := ctx.ShouldBindJSON(srv)
	if err == nil {
		ctx.JSON(http.StatusOK, srv.TagEdit())
		return
	}
	ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
	return
}

func (tc *TagsController) Build(rc *web.RouterCenter) {
	tagGroup := rc.RG.Group(consts.TagControllerGroup)
	tagGroup.Handle("POST", "/add", tc.tagAdd)
	tagGroup.Handle("DELETE", "/del/:id", tc.tagDel)
	tagGroup.Handle("GET", "/tag_list", tc.tagList)
	tagGroup.Handle("PUT", "/edit", tc.tagEdit)
}
