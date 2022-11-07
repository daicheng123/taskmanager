package ansible

//
//import (
//	"context"
//	"github.com/apenella/go-ansible/pkg/adhoc"
//	"github.com/apenella/go-ansible/pkg/execute"
//	"github.com/apenella/go-ansible/pkg/options"
//	"strings"
//	"time"
//)
//
//const (
//	StrictHostKeyChecking = "-o StrictHostKeyChecking=no"
//	SmartConnection       = "smart"
//	ForkNumber            = "10"
//	BecomeUser            = "root"
//)
//
//type AnsibleTask struct {
//	cmd             *adhoc.AnsibleAdhocCmd
//	OperateOvertime time.Duration
//	//callback ResultCallback
//}
//
//func NewAnsibleTask(options ...TaskOptions) *AnsibleTask {
//	task := &AnsibleTask{}
//	task.initAdhocCmd()
//	for _, o := range options {
//		o(task)
//	}
//	return task
//}
//
//func (at *AnsibleTask) initAdhocCmd() {
//	connOptions := &options.AnsibleConnectionOptions{
//		SSHCommonArgs: StrictHostKeyChecking,
//		Timeout:       10,
//		Connection:    SmartConnection,
//	}
//
//	privilegeOptions := &options.AnsiblePrivilegeEscalationOptions{
//		Become:     true,
//		BecomeUser: BecomeUser,
//	}
//
//	at.cmd = &adhoc.AnsibleAdhocCmd{
//		ConnectionOptions:          connOptions,
//		PrivilegeEscalationOptions: privilegeOptions,
//	}
//}
//
//func (at *AnsibleTask) CompleteTask(module, args string, callBack ResultCallback, inventorys ...string) *AnsibleTask {
//	at.cmd.Options = &adhoc.AnsibleAdhocOptions{
//		Forks:      ForkNumber,
//		Inventory:  strings.Join(inventorys, ","),
//		ModuleName: module,
//		Args:       args,
//	}
//
//	at.cmd.Exec = execute.NewDefaultExecute(
//		execute.WithWrite(callBack),
//		execute.WithWriteError(callBack))
//	return at
//}
//
////func()
//// TransferFile 传输脚本文件
//func (at *AnsibleTask) transferFile(ctx context.Context, scriptPath string) {
//
//}
//
//func (at *AnsibleTask) RunAdhoc(ctx context.Context) {
//	// 异步任务的taskid未获取
//
//}
//
//func (at *AnsibleTask) run(ctx context.Context) {
//	//at.cmd.
//	//
//}
//
//type TaskOptions func(task *AnsibleTask)
//
////func WithOvertime(duration uint) TaskOptions {
////	return func(task *AnsibleTask) {
////		task.OperateOvertime = time.Duration(duration) * time.Second
////	}
////}
