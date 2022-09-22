package email

import (
	"errors"
	"strings"
)

//const (
//	HtmlMail uint8 = iota + 1
//	TextMail
//
//	MailSubject = "【运维任务管理平台】"
//)
//
//type MailBuilder struct {
//	mail      *gomail.Message
//	receivers []string // 消息推送人
//	subject   string   // 消息主题
//	mailType  uint8
//	msgBody   string // 消息主体
//}
//
//func (mb *MailBuilder) SetMsgBody(mailType uint8, msgBody string) *MailBuilder {
//	mb.mailType = mailType
//	mb.msgBody = msgBody
//	return mb
//}
//
//func (mb *MailBuilder) Build() *MailService {
//	ms := &MailService{
//		receivers: strings.Join(mb.receivers, ","),
//		subject:   mb.subject,
//		mail:      mb.mail,
//	}
//	if mb.msgBody != "" {
//		ms.msgBody = mb.msgBody
//	}
//
//	switch mb.mailType {
//	case HtmlMail:
//		ms.mailType = "text/html"
//	//case consts.TextMail:
//	//	ms.mailType = "text/plain"
//	default:
//		ms.mailType = "text/plain"
//	}
//	return ms
//}
//
//func NewMailBuilder(receivers []string, subject string) *MailBuilder {
//	if len(receivers) == 0 || subject == "" {
//		return nil
//	}
//	mail := gomail.NewMessage(gomail.SetCharset("UTF-8"))
//	return &MailBuilder{
//		receivers: receivers,
//		subject:   subject,
//		mail:      mail,
//	}
//}
//
//type MailService struct {
//	receivers string // 消息推送人
//	subject   string // 消息主题
//	msgBody   string // 消息主体
//	mailType  string // html / plain
//	mail      *gomail.Message
//}
//
//func (ms *MailService) Builder(receivers []string, subject string) *MailBuilder {
//	return NewMailBuilder(receivers, subject)
//}
//
//func (ms *MailService) Sender() error {
//	ms.mail.SetHeader("Subject", ms.subject)
//	ms.mail.SetHeader("From", conf.GetMailUserName())
//	ms.mail.SetHeader("To", ms.receivers)
//	ms.mail.SetBody(ms.mailType, ms.msgBody)
//
//	dialer := gomail.NewDialer(
//		conf.GetMailHost(),
//		conf.GetMailPort(),
//		conf.GetMailUserName(),
//		conf.GetMailPwd())
//	return dialer.DialAndSend(ms.mail)
//}

var (
	// ErrChanNotOpen 邮件队列未开启
	ErrChanNotOpen = errors.New("邮件队列未开启")
	// ErrNoActiveDriver 无可用邮件发送服务
	ErrNoActiveDriver = errors.New("无可用邮件发送服务")
)

// MailDriver 邮件发送驱动
type MailDriver interface {
	// Close 关闭驱动
	Close()
	// Send 发送邮件
	Send(to, title, body string) error
}

// Send 发送邮件
func Send(receivers []string, title, body string) error {
	if len(receivers) == 0 {
		return errors.New("收件人不能唯空")
	}
	to := strings.Join(receivers, "")

	Lock.RLock()
	defer Lock.RUnlock()

	if MailClient == nil {
		return ErrNoActiveDriver
	}

	return MailClient.Send(to, title, body)
}
