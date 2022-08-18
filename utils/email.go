package utils

import (
	"gopkg.in/gomail.v2"
	"strings"
	"taskmanager/internal/conf"
	"taskmanager/internal/consts"
)

type MailBuilder struct {
	mail      *gomail.Message
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
		receivers: strings.Join(mb.receivers, ","),
		subject:   mb.subject,
		mail:      mb.mail,
	}
	if mb.msgBody != "" {
		ms.msgBody = mb.msgBody
	}

	switch mb.mailType {
	case consts.HtmlMail:
		ms.mailType = "text/html"
	//case consts.TextMail:
	//	ms.mailType = "text/plain"
	default:
		ms.mailType = "text/plain"
	}
	return ms
}

func NewMailBuilder(receivers []string, subject string) *MailBuilder {
	if len(receivers) == 0 || subject == "" {
		return nil
	}
	mail := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	return &MailBuilder{
		receivers: receivers,
		subject:   subject,
		mail:      mail,
	}
}

type MailService struct {
	receivers string // 消息推送人
	subject   string // 消息主题
	msgBody   string // 消息主体
	mailType  string // html / plain
	mail      *gomail.Message
}

func (ms *MailService) Builder(receivers []string, subject string) *MailBuilder {
	return NewMailBuilder(receivers, subject)
}

func (ms *MailService) Sender() error {
	ms.mail.SetHeader("Subject", ms.subject)
	ms.mail.SetHeader("From", conf.GetMailUserName())
	ms.mail.SetHeader("To", ms.receivers)
	ms.mail.SetBody(ms.mailType, ms.msgBody)

	dialer := gomail.NewDialer(
		conf.GetMailHost(),
		conf.GetMailPort(),
		conf.GetMailUserName(),
		conf.GetMailPwd())
	return dialer.DialAndSend(ms.mail)
}
