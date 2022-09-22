package serializer

import (
	"sync"
)

type Resolve func(rsp interface{})

type Reject func(res *Response) // 异常处理函数

type PromiseFunc func(resolve Resolve, reject Reject)

type Promise struct {
	f       PromiseFunc // 业务处理函数
	resolve Resolve     // resolve 函数
	reject  Reject      // reject 函数
	wg      sync.WaitGroup
}

//Then  业务处理函数状态为resolve时调用 链式调用
func (p *Promise) Then(resolve Resolve) *Promise {
	p.resolve = resolve
	return p
}

//Catch  业务处理函数状态为reject时调用 链式调用
func (p *Promise) Catch(reject Reject) *Promise {
	p.reject = reject
	return p
}

func (p *Promise) Done() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.f(p.resolve, p.reject)
	}()
	p.wg.Wait()
}

func NewPromise(f PromiseFunc) *Promise {
	return &Promise{
		f: f,
	}
}

//func test() {
//	NewPromise(func(resolve Resolve, reject Reject) {
//		time.Sleep(time.Second * 1)   // 模拟业务处理
//		if time.Now().Unix()%2 == 0 { // 模拟业务处理成功或失败
//			resolve("ok")
//		} else {
//			reject(errors.New("my error"))
//		}
//	}).Then(func(rsp interface{}) {
//		fmt.Println(rsp)
//	}).Catch(func(err error) {
//		log.Println(err)
//	}).Done()
//}
