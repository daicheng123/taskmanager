package ansible

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"io"
	"os"
	"strings"
	"taskmanager/internal/consts"
	"taskmanager/pkg/ansible/templates"
	"taskmanager/utils"
	"text/template"
)

var (
	OperationTimeoutError = errors.New("operation timeout")
)

type OperatorConfig struct {
	taskName  string
	Executor  uint
	sshUser   string
	tpls      []string
	filePath  string
	overtime  int
	inventory string
	ResultCallback
}

type OperatorOptions func(cfg *OperatorConfig)

func WithOvertime(duration int) OperatorOptions {
	return func(cfg *OperatorConfig) {
		cfg.overtime = duration
	}
}

func WithTaskName(name string) OperatorOptions {
	return func(cfg *OperatorConfig) {
		cfg.taskName = name
	}
}

func WithTaskUser(userName string) OperatorOptions {
	return func(cfg *OperatorConfig) {
		cfg.sshUser = userName
	}
}

func WithTaskExecutor(executorId uint) OperatorOptions {
	return func(cfg *OperatorConfig) {
		cfg.Executor = executorId
	}
}

func (c *OperatorConfig) AddInventory(inventorys ...string) {
	c.inventory = strings.Join(inventorys, ",")
}

func (c *OperatorConfig) SetCallBack(fn func(tn string) ResultCallback) {
	c.ResultCallback = fn(c.taskName)
}

func NewConfig(options ...OperatorOptions) *OperatorConfig {
	c := new(OperatorConfig)

	for _, option := range options {
		option(c)
	}
	if c.overtime <= 0 {
		c.overtime = 120
	}

	if len(c.taskName) == 0 {
		c.taskName = utils.NewUuid()
	}
	return c
}

func (c *OperatorConfig) generateScript(content, path string) error {
	return utils.WriteFile(path, content)
}

func (c *OperatorConfig) delete() {
	os.Remove(c.filePath)
	for _, tpl := range c.tpls {
		os.Remove(tpl)
	}
}

func (c *OperatorConfig) RenderPlayBook(scriptContent string, scriptType consts.ScriptType) error {
	var (
		prefix      = utils.BuilderStr("/tmp/", c.taskName, c.inventory)
		suffix      string
		interpreter string
	)
	if scriptType == consts.Python {
		suffix = ".py"
		interpreter = "python3 "
	} else if scriptType == consts.Shell {
		suffix = ".sh"
		interpreter = "bash "
	}
	c.filePath = utils.BuilderStr(prefix, suffix)
	if err := c.generateScript(scriptContent, c.filePath); err != nil {
		return err
	}
	cmd := utils.BuilderStr(interpreter, c.filePath)

	t := struct {
		Src     string `json:"src"`
		Dest    string `json:"dest"`
		Command string `json:"command"`
	}{c.filePath, c.filePath, cmd}

	tpl, err := template.New("operation").Parse(templates.OperationTpl)
	if err != nil {
		return err
	}

	buff := &bytes.Buffer{}
	err = tpl.Execute(buff, t)
	if err != nil {
		return err
	}
	tplPath := utils.BuilderStr(prefix, ".yaml")

	err = utils.WriteFile(tplPath, buff.String())
	if err == nil {
		c.tpls = append(c.tpls, tplPath)
	}
	return err
}

type OperationTask struct {
	cmd *playbook.AnsiblePlaybookCmd
	*OperatorConfig
}

func NewOperationTask(c *OperatorConfig) *OperationTask {
	task := &OperationTask{
		OperatorConfig: c,
	}
	return task
}

func (t *OperationTask) Execute(ctx context.Context) error {
	if t.OperatorConfig == nil {
		return errors.New("ansible config can not be nil")
	}
	connOptions := &options.AnsibleConnectionOptions{
		SSHCommonArgs: StrictHostKeyChecking,
		Timeout:       t.OperatorConfig.overtime,
		User:          t.OperatorConfig.sshUser,
		Connection:    SmartConnection,
	}
	fmt.Printf("executorId is: %#v", t.OperatorConfig.inventory)
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: t.OperatorConfig.inventory + ",",
		Forks:     ForkNumber,
	}

	privilegeOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become:     true,
		BecomeUser: BecomeUser,
	}

	buff := new(bytes.Buffer)
	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithEnvVar(AnsibleRetryFilesEnabled, "0"),
			execute.WithWrite(io.Writer(buff)),
			execute.WithWriteError(io.Writer(buff)),
		),
	)

	t.cmd = &playbook.AnsiblePlaybookCmd{
		Playbooks:                  t.OperatorConfig.tpls,
		Exec:                       executorTimeMeasurement,
		ConnectionOptions:          connOptions,
		Options:                    ansiblePlaybookOptions,
		PrivilegeEscalationOptions: privilegeOptions,
		StdoutCallback:             "json",
	}

	ret := make(chan error, 1)

	go utils.RunSafeWithMsg(func() {
		ret <- t.cmd.Run(ctx)
	}, "ansible operation task failed to run")

	select {
	case <-ctx.Done():
		return OperationTimeoutError
	case <-ret:
		t.OperatorConfig.ResultCallback.ParseResult(ctx, buff)

		go utils.RunSafeWithMsg(func() {
			t.OperatorConfig.delete()
		}, "failed to delete script")
		return nil
	}
}
