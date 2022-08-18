package utils

type Result struct {
	Result interface{}
	Err    error
}

type Object interface {
	bool
}

func NewResult(result interface{}, err error) *Result {
	return &Result{Result: result, Err: err}
}

func (sr *Result) Unwrap() string {
	if sr.Err != nil {
		panic(sr.Err)
	}
	return sr.Result.(string)
}

//UnwrapOr  如果获取出错，则返回一个默认值
func (sr *Result) UnwrapOr(str string) string {
	if sr.Err != nil {
		return str
	}
	return sr.Result.(string)
}

func (sr *Result) UnwrapOrElse(f func(err error)) interface{} {
	if sr.Err != nil {
		f(sr.Err)
	}
	return sr.Result
}
