package utils

import (
	"gopkg.in/gomail.v2"
	"strings"
	"taskmanager/internal/conf"
)

const (
	HtmlMail uint8 = 1
	TextMail uint8 = 2
)

type MailBuilder struct {
	receivers []string // 消息推送人
	subject   string   // 消息主题
	mailType  uint8
	msgBody   string // 消息主体
}

func (mb *MailBuilder) SetMsgBody(mailType uint8, msgBody string) *MailBuilder {
	mb.mailType = mailType
	mb.msgBody = msgBody
	return mb
}

func (mb *MailBuilder) Build() *MailService {
	ms := &MailService{
		senders: strings.Join(mb.receivers, ","),
		subject: mb.subject,
	}
	if mb.msgBody != "" {
		ms.msgBody = mb.msgBody
	}

	switch mb.mailType {
	case HtmlMail:
		ms.mailType = "text/html"
	case TextMail:
		ms.mailType = "text/plain"
	default:
		ms.mailType = "text/plain"
	}
	return ms
}

func NewMailBuilder(receivers []string, subject string) *MailBuilder {
	if len(receivers) == 0 || subject == "" {
		return nil
	}

	return &MailBuilder{
		receivers: receivers,
		subject:   subject,
	}
}

type MailService struct {
	senders  string // 消息推送人
	subject  string // 消息主题
	msgBody  string // 消息主体
	mailType string // html / text
	mail     *gomail.Message
}

func (ms *MailService) Builder(senders []string, subject string) *MailBuilder {
	return NewMailBuilder(senders, subject)
}

func (ms *MailService) Sender() error {
	ms.mail.SetHeader("Subject", ms.subject)
	ms.mail.SetHeader("From", ms.senders)
	ms.mail.SetBody(ms.mailType, ms.msgBody)

	dialer := gomail.NewDialer(
		conf.GetMailHost(),
		conf.GetMailPort(),
		conf.GetMailUserName(),
		conf.GetMailPwd())

	return dialer.DialAndSend(ms.mail)
}
