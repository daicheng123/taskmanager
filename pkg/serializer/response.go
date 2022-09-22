package serializer

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg"`
	Error   string      `json:"error,omitempty"`
}
