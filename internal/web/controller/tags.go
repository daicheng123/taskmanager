package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/consts"
	"taskmanager/internal/service/admin"
	"taskmanager/internal/web"
	"taskmanager/internal/web/utils"
)

/*
	TagsController 脚本标签管理
*/

type TagsController struct {
}

func NewTagsController() *TagsController {
	return &TagsController{}
}

func (tc *TagsController) TagAdd(ctx *gin.Context) {
	srv := &admin.TagsService{}
	err := ctx.ShouldBindJSON(srv)
	fmt.Printf("data: %+v\n", srv)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, srv.TagSave())
}

func (tc *TagsController) TagList(ctx *gin.Context) {
	service := &admin.TagsService{}
	err := ctx.ShouldBindQuery(service)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, service.ListByPager())
}

func (tc *TagsController) Build(rc *web.RouterCenter) {
	tagGroup := rc.RG.Group(consts.TagControllerGroup)
	tagGroup.Handle("POST", "/add", tc.TagAdd)
	tagGroup.Handle("DELETE", "/del/:id", tc.TagAdd)
	tagGroup.Handle("GET", "/tag_list", tc.TagList)
	tagGroup.Handle("GET", "/detail/:id", tc.TagAdd)
	tagGroup.Handle("PATCH", "/edit/", tc.TagAdd)
}
