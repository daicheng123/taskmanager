package core

import (
	"github.com/gorilla/websocket"
	"net/http"
	"taskmanager/pkg/logger"
	"time"
)

var UpGrader websocket.Upgrader

func init() {
	UpGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 可以检验是否允许跨域的来源ip
			return true
		},
	}
}

type WebSocketClient struct {
	connect   *websocket.Conn
	readChan  chan WsMessage
	closeChan chan struct{}
}

func NewWebSocketClient(conn *websocket.Conn) *WebSocketClient {
	client := &WebSocketClient{
		connect:   conn,
		readChan:  make(chan WsMessage),
		closeChan: make(chan struct{}),
	}
	return client
}

func (wsc *WebSocketClient) Send(msg WsMessage) {
	wsc.readChan <- msg
}

func (wsc *WebSocketClient) HandlerLoop() {
LOOP:
	for {
		select {
		case msg := <-wsc.readChan:
			if err := wsc.connect.WriteMessage(websocket.TextMessage, msg.Render()); err != nil {
				wsc.closeChan <- struct{}{}
				ClientMap.Remove(wsc)
			}
		case <-wsc.closeChan:
			logger.Warning("client %s has closed", wsc.connect.RemoteAddr())
			break LOOP
		}
	}
}

func (wsc *WebSocketClient) PingLoop(dur time.Duration) {
	ticker := time.NewTicker(dur)
	for {
		select {
		case <-ticker.C:
			err := wsc.connect.WriteMessage(websocket.TextMessage, []byte("ping"))
			if err != nil {
				ClientMap.Remove(wsc)
				return
			}
		}
	}
}
