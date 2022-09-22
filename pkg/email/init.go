package email

import (
	"sync"
	"taskmanager/internal/conf"
)

func init() {
	InitDriver()
}

var (
	MailClient MailDriver
	Lock       sync.RWMutex
)

func InitDriver() {
	Lock.Lock()
	defer Lock.Unlock()
	if MailClient != nil {
		MailClient.Close()
	}

	client := NewSMTPClient(SMTPConfig{
		Name:       "自动化任务平台",
		Address:    conf.GetMailHost(),
		ReplyTo:    conf.GetMailHost(),
		Host:       conf.GetMailHost(),
		Port:       conf.GetMailPort(),
		User:       conf.GetMailUserName(),
		Password:   conf.GetDbPassword(),
		Encryption: true,
		Keepalive:  30,
	})
	MailClient = client
}
