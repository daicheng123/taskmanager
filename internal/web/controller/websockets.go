package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskmanager/internal/web"
	"taskmanager/pkg/serializer"
	"taskmanager/pkg/websockets"
)

const (
	wsControllerGroup = "websockets"
)

type WebSocketController struct {
}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{}
}

func (wsc *WebSocketController) Connect(ctx *gin.Context) {
	//conn, err := core.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)

	err := websockets.Echo(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusOK, &serializer.Response{
			Code: http.StatusInternalServerError, Message: "websocket创建失败", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &serializer.Response{Message: "ok"})
}

func (wsc *WebSocketController) Build(rc *web.RouterCenter) {
	g := rc.Group(wsControllerGroup)
	g.Handle(http.MethodGet, "/core", wsc.Connect)
}
