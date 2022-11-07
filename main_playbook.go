package main

import (
	"bytes"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"io"
	"log"
	"taskmanager/pkg/ansible"
	"time"

	//"bytes"
	"context"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

func main() {

	//var err error
	//var res *results.AnsiblePlaybookJSONResults
	buff := new(bytes.Buffer)

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection:    "ssh",
		User:          "devops",
		SSHCommonArgs: "-o StrictHostKeyChecking=no",
		Timeout:       1,
	}

	ansiblePrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become:     true,
		BecomeUser: "root",
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: "172.31.176.42,172.31.176.4,",
		Forks:     "10",
	}
	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
			execute.WithWriteError(io.Writer(buff)),
		),
	)
	playbooksList := []string{"./site.yaml"}
	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  playbooksList,
		PrivilegeEscalationOptions: ansiblePrivilegeEscalationOptions,
		Exec:                       executorTimeMeasurement,
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		Options:                    ansiblePlaybookOptions,
		StdoutCallback:             "json",
	}
	options.AnsibleSetEnv("ANSIBLE_RETRY_FILES_ENABLED", "0")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*200)
	defer cancelFunc()
	err := playbook.Run(ctx)
	if err != nil {
		log.Println(err.Error())
	}
	//r.DoneChan <- struct{}{}
	//r.ParseResult()
	//fmt.Println(buff.String())
	//fmt.Println(buff.String())
	cb := ansible.NewOperationCallback(&ansible.OperationCallbackPayload{})

	cb.ParseResult(ctx, buff)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//
	//results.AnsiblePlaybookJSONResults{}
	//fmt.Println(res.String())
	//fmt.Println("Duration: ", executorTimeMeasurement.Duration())
}
