package utils

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"taskmanager/pkg/serializer"
	"time"
)

const (
	keyCmd = "if test ! -d  $HOME/.ssh; then mkdir $HOME/.ssh; fi && echo %s >> $HOME/.ssh/authorized_keys"
)

var (
	ErrGetSecretFileFailed = serializer.NewError(serializer.CodeHostSecretFileNotFound, "密钥文件不存在", nil)
	ErrCreateSessionFailed = serializer.NewError(serializer.CodeHostUnreachable, "创建主机会话失败", nil)
	ErrCommandExecFailed   = serializer.NewError(serializer.CodeHostCommandErr, "执行命令行错误", nil)
)

type SSHClient struct {
	UserName     string
	UserPassword string
	HostAddr     string
	HostPort     uint
}

func NewSSHClient(name, pwd, addr string, port uint) *SSHClient {
	return &SSHClient{
		UserName:     name,
		UserPassword: pwd,
		HostAddr:     addr,
		HostPort:     port,
	}
}
func (sc *SSHClient) NewSSHSession() (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(sc.UserPassword))

	clientConfig = &ssh.ClientConfig{
		User:            sc.UserName,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr = fmt.Sprintf("%s:%d", sc.HostAddr, sc.HostPort)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, ErrCreateSessionFailed.WithError(err)
	}

	if session, err = client.NewSession(); err != nil {
		return nil, ErrCreateSessionFailed.WithError(err)
	}
	return session, nil
}

func (sc *SSHClient) DistributeKey() (outStr string, err error) {
	f, err := ReadFile("/Users/daicheng/.ssh/id_rsa.pub")
	if err != nil {
		return "", ErrGetSecretFileFailed.WithError(err)
	}
	cmd := fmt.Sprintf(keyCmd, f)
	return sc.RemoteCommand(cmd)
}

func (sc *SSHClient) RemoteCommand(cmd string) (outStr string, err error) {
	session, err := sc.NewSSHSession()
	if err != nil {
		return "", err
	}
	var (
		stdOut = &bytes.Buffer{}
		stdErr = &bytes.Buffer{}
	)
	session.Stdout = stdOut
	session.Stderr = stdErr
	err = session.Run(cmd)
	if stdErr == nil || err == nil {
		return stdOut.String(), nil
	}

	errStr := ""
	if stdErr != nil {
		errStr = fmt.Sprintf("error: %s", err.Error())
	}
	if stdErr != nil {
		errStr = fmt.Sprintf("%s, stdError: %s", stdErr.String())
	}
	err = errors.New(errStr)
	return "", ErrCommandExecFailed.WithError(err)
}
