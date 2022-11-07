package serializer

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"msg"`
	Error   string      `json:"error,omitempty"`
}

//
//func HandleResponse(ctx *gin.Context, fn func() (interface{}, error)) {
//	resp := &Response{
//		Code:    http.StatusOK,
//		Message: "",
//	}
//
//	defer func() {
//		if p := recover(); p != nil {
//			resp.Message = fmt.Sprintf("%s", p)
//		}
//	}()
//
//	body, err := fn()
//	var (
//		appErr  AppError
//		code    int
//		message string
//	)
//	if err != nil {
//		if errors.As(err, &appErr) {
//			code = appErr.Code
//			message = appErr.Message
//			err = appErr.RawError
//		}
//
//		if gin.Mode() != gin.ReleaseMode {
//			resp.Error = err.Error()
//		}
//	}
//
//	resp.Data = body
//}
