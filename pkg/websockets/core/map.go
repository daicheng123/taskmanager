package core

import (
	"sync"
	"taskmanager/pkg/logger"
	"time"
)

type WsClientMap struct {
	data sync.Map
}

//ClientMap 外部公共使用
var ClientMap *WsClientMap

func init() {
	ClientMap = &WsClientMap{}
}

func (wm *WsClientMap) Store(connCli *WebSocketClient) {
	wm.data.Store(connCli.connect.RemoteAddr(), connCli)
	go connCli.HandlerLoop()
	go connCli.PingLoop(time.Second * 5)
}

func (wm *WsClientMap) Remove(connCli *WebSocketClient) {
	wm.data.Delete(connCli.connect.RemoteAddr())
}

func (wm *WsClientMap) SendAll(msg WsMessage) {
	wm.data.Range(func(key, value any) bool {
		c := value.(*WebSocketClient)
		c.Send(msg)
		return true
	})
}

func (wm *WsClientMap) SendByRemoteAddr(addr string, msg WsMessage) {
	if value, ok := wm.data.Load(addr); ok {
		value.(*WebSocketClient).Send(msg)
		return
	}
	logger.Warning("can not found client %s in local store", addr)
}
