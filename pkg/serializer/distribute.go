package serializer

type DistributeRes struct {
	IsHostAdd    bool   `json:"isHostAdd"`
	SecretStatus uint   `json:"secretStatus"`
	ErrMsg       string `json:"errMsg"`
}

//type DistributeFunc func(*DistributeRes)
//
//type DistributeFuncs []DistributeFunc
//
//func (drs DistributeFuncs) apply() {
//	for _, f := range drs {
//
//	}
//}
//
//func WithErrMsg(errMsg string) DistributeFunc {
//	return func(dr *DistributeRes) {
//		dr.ErrMsg = errMsg
//	}
//}
//
//func WithSecretStatus(secretStatus uint) DistributeFunc {
//	return func(dr *DistributeRes) {
//		dr.SecretStatus = secretStatus
//	}
//}
//
//func WithIsHostAdd(isAdd bool) DistributeFunc {
//	return func(dr *DistributeRes) {
//		dr.IsHostAdd = isAdd
//	}
//}
//
//func (dr *DistributeRes) Apply(fns ...DistributeFunc) {
//	for _, fn := range fns {
//		fn(dr)
//	}
//}
//
//func (dr *DistributeRes) Builder() *Response {
//	var rsp = &Response{}
//	rsp.Data = dr
//	return rsp
//}

//func BuildDistributeRes(executor *models.Executor, secretStatus uint) Response{
//	var dr = new(DistributeRes)
//	dr.SecretStatus = secretStatus
//	dr.ErrMsg =
//}
