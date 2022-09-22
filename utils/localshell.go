package utils

import (
	"bytes"
	"os/exec"
)

//ExecuteShell  https://chengchaos.github.io/2020/08/10/golang-execute-cmd.html
func ExecuteShell(cmdStr string) (out string, err error) {
	cmd := exec.Command("sh", "-c", cmdStr)
	var (
		stdOut bytes.Buffer
	)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdOut
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return stdOut.String(), err
}
