package core

import (
	ws "github.com/gorilla/websocket"
	"sync"
	"time"
)

type ClientMapStruct struct {
	data sync.Map //  key 是客户端IP; value 就是 WsClient连接对象
}

//ClientMap 外部公共使用
var ClientMap *ClientMapStruct

func init() {
	ClientMap = &ClientMapStruct{}
}

func (cms *ClientMapStruct) Store(conn *ws.Conn) {
	wsCil := NewWsClient(conn)
	cms.data.Store(conn.RemoteAddr().String(), wsCil)
	// 与客户端维持心跳
	go wsCil.Ping(time.Second * 5)
	// 读取客户端数据
	go wsCil.ReadLoop()
	// 独立响应客户端的循环
	//go wsCil.WriteLoop()
	// 处理客户端数据
	go wsCil.HandlerLoop()
}

func (cms *ClientMapStruct) Remove(conn *ws.Conn) {
	cms.data.Delete(conn.RemoteAddr().String())
}
