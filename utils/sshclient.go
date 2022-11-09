package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"taskmanager/internal/conf"
	"taskmanager/pkg/logger"
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
	// ErrTransferFileFailed  = serializer.NewError(serializer.CodeHostCommandErr, "执行命令行错误", nil)
)

type SSHClient struct {
	UserName     string
	UserPassword string
	HostAddr     string
	HostPort     uint
}

func NewSsh(name, pwd, addr string, port uint) *SSHClient {
	return &SSHClient{
		UserName:     name,
		UserPassword: pwd,
		HostAddr:     addr,
		HostPort:     port,
	}
}
func (sc *SSHClient) NewSSHClient() (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		//session      *ssh.Session
		err error
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
		logger.Error("创建ssh会话失败, err: [%s]", err.Error())
		return nil, ErrCreateSessionFailed.WithError(err)
	}
	return client, nil
}

func (sc *SSHClient) DistributeKey() (outStr string, err error) {
	cmd := fmt.Sprintf(keyCmd, conf.GetSecrets())
	return sc.RemoteCommand(cmd)
}

func (sc *SSHClient) RemoteCommand(cmd string) (outStr string, err error) {
	client, err := sc.NewSSHClient()
	if err != nil {
		return "", err
	}
	var (
		stdOut = &bytes.Buffer{}
		stdErr = &bytes.Buffer{}
	)
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	session.Stdout = stdOut
	session.Stderr = stdErr
	err = session.Run(cmd)
	if stdErr == nil || err == nil {
		return stdOut.String(), nil
	}

	errStr := ""
	if err != nil {
		errStr = fmt.Sprintf("error: %s", err.Error())
	}
	if stdErr != nil {
		errStr = fmt.Sprintf("%s, stdError: %s", errStr, stdErr.String())
	}
	err = errors.New(errStr)
	return "", ErrCommandExecFailed.WithError(err)
}

func (sc *SSHClient) GetSftpClient() (*sftp.Client, error) {
	client, err := sc.NewSSHClient()
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(
		client,
		sftp.MaxConcurrentRequestsPerFile(20),
		sftp.MaxPacket(1000))

	return sftpClient, err
}

func (sc *SSHClient) TransferFile(src, dest string) error {
	client, err := sc.GetSftpClient()
	if err != nil {
		return err
	}
	defer client.Close()
	file, err := client.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadFile(src)
	file.Write(bytes)
	return nil
}
