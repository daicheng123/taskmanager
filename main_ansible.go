package main

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/adhoc"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/stdoutcallback"
	"log"

	//"github.com/spf13/cobra"
	//"github.com/spf13/viper"
	"io"
	"os"
)

//func main() {
//	ansibleConnectionOptions := &options.AnsibleConnectionOptions{
//		Connection: "local",
//	}
//
//	ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
//		Inventory:  " 127.0.0.1,",
//		ModuleName: "command",
//		Args:       "ping 127.0.0.1 -c 2",
//	}
//
//	adhoc := &adhoc.AnsibleAdhocCmd{
//		Pattern:           "all",
//		Options:           ansibleAdhocOptions,
//		ConnectionOptions: ansibleConnectionOptions,
//		StdoutCallback:    "yaml",
//	}
//
//	fmt.Println("Command: ", adhoc.String())
//
//	err := adhoc.Run(context.TODO())
//	if err != nil {
//		panic(err)
//	}
//}

type MyExecutor struct {
	Prefix string
}

func (e *MyExecutor) Options(options ...execute.ExecuteOptions) {
	// apply all options to the executor
	for _, opt := range options {
		opt(e)
	}
}

func WithPrefix(prefix string) execute.ExecuteOptions {
	return func(e execute.Executor) {
		e.(*MyExecutor).Prefix = prefix
	}
}

func (e *MyExecutor) Execute(ctx context.Context, command []string, resultsFunc stdoutcallback.StdoutCallbackResultsFunc, options ...execute.ExecuteOptions) error {

	// apply all options to the executor
	for _, opt := range options {
		opt(e)
	}

	fmt.Println(fmt.Sprintf("%s %s\n", e.Prefix, "I am MyExecutor and I am doing nothing"))
	//resultsFunc()
	return nil
}

const (
	ErrorPathUnknown = "path does not exist"

	flagAnsibleRoot   = "ansible-root-directory"
	flagInventoryFile = "ansible-inventory-file"
	flagPlaybookFile  = "ansible-playbook-file"
)

type Config struct {
	RootDir       string
	InventoryFile string
	PlaybookFile  string

	stdout io.Writer
	stderr io.Writer

	// ExtraVars is a key value map that will be passed to Ansible
	// as extra variable using --extra-vars.
	// The corresponding keys are defined as constants in the `vars.go` file of this package.
	// This map gets filled by the Executor before the execution of Ansible.
	// It contains all the required information to run the Ansible roles and tasks properly.
	ExtraVars map[string]interface{}
}

func NewConfig() Config {
	return Config{
		ExtraVars: map[string]interface{}{},
	}
}

func (c *Config) AddExtraVar(key string, value interface{}) {
	c.ExtraVars[key] = value
}

//func (c *Config) AddToViper(v *viper.Viper) {
//	_ = v.UnmarshalKey(flagAnsibleRoot, &c.RootDir)
//	_ = v.UnmarshalKey(flagInventoryFile, &c.InventoryFile)
//	_ = v.UnmarshalKey(flagPlaybookFile, &c.PlaybookFile)
//}
//
//func (c *Config) AddToPersistentCommand(cmd *cobra.Command) {
//	cmd.Flags().StringVar(&c.RootDir, flagAnsibleRoot, "", "Root directory of Ansible")
//	cmd.Flags().StringVar(&c.InventoryFile, flagInventoryFile, "", "Inventory file used by Ansible")
//	cmd.Flags().StringVar(&c.PlaybookFile, flagPlaybookFile, "", "Playbook file used by Ansible")
//
//	_ = viper.BindPFlag(flagAnsibleRoot, cmd.Flags().Lookup(flagAnsibleRoot))
//	_ = viper.BindPFlag(flagInventoryFile, cmd.Flags().Lookup(flagInventoryFile))
//	_ = viper.BindPFlag(flagPlaybookFile, cmd.Flags().Lookup(flagPlaybookFile))
//}

//func applyRootToFiles(root string, file *string) {
//	if !path.IsAbs(*file) {
//		*file = path.Join(root, *file)
//	}
//}

func Run(c *Config) error {
	//applyRootToFiles(c.RootDir, &c.PlaybookFile)
	//applyRootToFiles(c.RootDir, &c.InventoryFile)

	ansibleAdhocConnectionOptions := &options.AnsibleConnectionOptions{
		User:          "devops",
		SSHCommonArgs: "-o StrictHostKeyChecking=no",
	}

	//ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
	//	Inventory: c.InventoryFile,
	//	ExtraVars: c.ExtraVars,
	//}
	ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
		Inventory:  "172.31.176.42 ,",
		ModuleName: "command",
		Args:       "echo hello",
	}

	ansibleAdhocPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become:     true,
		BecomeUser: "root",
	}

	ahc := &adhoc.AnsibleAdhocCmd{
		//Playbooks:                  []string{c.PlaybookFile},
		Pattern:                    "all",
		ConnectionOptions:          ansibleAdhocConnectionOptions,
		PrivilegeEscalationOptions: ansibleAdhocPrivilegeEscalationOptions,
		Options:                    ansibleAdhocOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithShowDuration(),
			execute.WithWrite(c.stdout),
			execute.WithWriteError(c.stderr),
			//execute.WithShowDuration(),
			//Exec: MyExecutor{},
		),
		StdoutCallback: "debug",
	}

	err := ahc.Run(context.TODO())
	fmt.Println(ahc.String())
	if err != nil {
		return err
	}
	return nil
}

//func (c *Config) CopyRootDirectory(directory string) error {
//	if _, err := os.Stat(directory); os.IsNotExist(err) {
//		return errors.New(ErrorPathUnknown)
//	}
//
//	err := copy.Copy(c.RootDir, directory)
//	if err != nil {
//		return err
//	}
//	c.RootDir = directory
//	return nil
//}

func (c *Config) SetStdout(stdout *os.File) {
	c.stdout = stdout
}

func (c *Config) SetStderr(stderr *os.File) {
	c.stderr = stderr
}

func (c *Config) SetOutputs(stdout, stderr *os.File) {
	c.stdout = stdout
	c.stderr = stderr
}

func main() {
	file, err := os.OpenFile("./test.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	c := NewConfig()
	c.SetOutputs(file, file)

	err = Run(&c)
	fmt.Println(err)
}
