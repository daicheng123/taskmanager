package models

import "net/http"

type WebsocketServer struct {
}

func NewWebsocketServer() *WebsocketServer {
	return &WebsocketServer{}
}

func SocketHandler(hub *SocketHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWebsocket(hub, w, r)
	}
}
