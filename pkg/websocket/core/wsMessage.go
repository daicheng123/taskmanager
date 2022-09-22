package core

type WsMessage struct {
	MessageType int    // 这个基本用不到
	MessageDate []byte // 消息体，一般是 json 格式
}

func NewWsMessage(messageType int, messageDate []byte) *WsMessage {
	return &WsMessage{MessageType: messageType, MessageDate: messageDate}
}
