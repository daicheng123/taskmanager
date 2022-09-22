package core

import (
	ws "github.com/gorilla/websocket"
	"log"
	"taskmanager/internal/models"
	"time"
)

type WsClient struct {
	conn *ws.Conn
	// 读取客户端数据的channel
	readChan chan *WsMessage
	// 停止channel
	closeChan chan struct{}
	// 写队列
	resultChan chan *models.WebsocketRes
}

func NewWsClient(conn *ws.Conn) *WsClient {
	return &WsClient{
		conn:      conn,
		closeChan: make(chan struct{}),
	}
}

func (wc *WsClient) Ping(dur time.Duration) {
	ticker := time.NewTicker(dur)
LOOP:
	for {
		select {
		case <-ticker.C:
			err := wc.conn.WriteMessage(ws.TextMessage, []byte("ping"))
			if err != nil {
				ClientMap.Remove(wc.conn)
				break LOOP
			}
		}
	}
}

func (wc *WsClient) ReadLoop() {
	for {
		t, data, err := wc.conn.ReadMessage()
		if err != nil {
			// 关闭该ws 连接
			wc.conn.Close()
			ClientMap.Remove(wc.conn)
			wc.closeChan <- struct{}{} // 出错通知HandlerLoop关闭
		}
		wc.readChan <- NewWsMessage(t, data)
	}
}

func (wc *WsClient) HandlerLoop() {
LOOP:
	for {
		select {
		case _ = <-wc.readChan:
			return
		case <-wc.closeChan:
			log.Println("has closed")
			break LOOP
		}
	}
}

//func (wc *WsClient) WriteLoop() {
//LOOP:
//	for {
//		select {
//		case msg := <-wc.resultChan:
//			if err := ws.Conn.WriteMessage(ws.TextMessage, msg); err != nil {
//				wc.conn.Close()
//				ClientMap.Remove(wc.conn)
//				wc.closeChan <- struct{}{}
//				break LOOP
//			}
//		}
//	}
//}
