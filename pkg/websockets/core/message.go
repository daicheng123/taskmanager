package core

type WsMessage interface {
	Render() []byte
}
