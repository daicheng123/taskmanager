package websockets

import (
	"net/http"
	"taskmanager/pkg/websockets/core"
)

// Echo 将升级为websocket后的客户端存入到本地内存的 ClientMap 中
func Echo(w http.ResponseWriter, req *http.Request) (err error) {
	conn, err := core.UpGrader.Upgrade(w, req, nil)
	if err != nil {
		return err
	} else {
		wsCli := core.NewWebSocketClient(conn)
		core.ClientMap.Store(wsCli) // 如果正常将 ws 连接对象保存至map中
	}
	return
}

