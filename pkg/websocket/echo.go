package websocket

import (
	"log"
	"net/http"
	"taskmanager/pkg/websocket/core"
)

// Echo 将升级为websocket后的客户端存入到本地内存的 ClientMap 中
func Echo(w http.ResponseWriter, req *http.Request) {
	//升级为ws
	client, err := core.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
	} else {
		core.ClientMap.Store(client) // 如果正常将 ws 连接对象保存至map中
	}
}
