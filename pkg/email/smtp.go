package email

import (
	gomail "github.com/go-mail/mail"
	"taskmanager/pkg/logger"
	"time"
)

type SMTP struct {
	Config SMTPConfig
	ch     chan *gomail.Message
	chOpen bool
}

type SMTPConfig struct {
	Name       string // 发送者名
	Address    string // 发送者地址
	ReplyTo    string // 回复地址
	Host       string // 邮件服务器主机名
	Port       int    // 邮件服务器端口
	User       string // 用户名
	Password   string // 密码
	Encryption bool   // 是否启用加密
	Keepalive  int    // SMTP 连接保留时长
}

// NewSMTPClient 新建SMTP发送队列
func NewSMTPClient(config SMTPConfig) *SMTP {
	client := &SMTP{
		Config: config,
		ch:     make(chan *gomail.Message, 30),
		chOpen: false,
	}

	client.Init()

	return client
}

func (client *SMTP) Send(to, title, body string) error {
	if !client.chOpen {
		return ErrChanNotOpen
	}
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", client.Config.Address, client.Config.Name)
	msg.SetAddressHeader("Reply-To", client.Config.ReplyTo, client.Config.Name)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", title)
	msg.SetHeader("text/html", body)
	client.ch <- msg
	return nil
}

// Close 关闭发送队列
func (client *SMTP) Close() {
	if client.ch != nil {
		close(client.ch)
	}
}

// Init 初始化发送队列
func (client *SMTP) Init() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				client.chOpen = false
				logger.Error("邮件发送队列出现异常, %s ,10 秒后重置", err)
				time.Sleep(time.Duration(10) * time.Second)
				client.Init()
			}
		}()

		d := gomail.NewDialer(client.Config.Host, client.Config.Port, client.Config.User, client.Config.Password)
		d.Timeout = time.Duration(client.Config.Keepalive+5) * time.Second
		client.chOpen = true
		// 是否启用 SSL
		d.SSL = false
		if client.Config.Encryption {
			d.SSL = true
		}
		d.StartTLSPolicy = gomail.OpportunisticStartTLS

		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case message, ok := <-client.ch:
				if !ok {
					logger.Debug("邮件队列关闭")
					client.chOpen = false
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				if err := gomail.Send(s, message); err != nil {
					logger.Error("邮件发送失败, %s", err)
				} else {
					logger.Debug("邮件已发送")
				}
			// 长时间没有新邮件，则关闭SMTP连接
			case <-time.After(time.Duration(client.Config.Keepalive) * time.Second):
				if open {
					if err := s.Close(); err != nil {
						logger.Warning("无法关闭 SMTP 连接 %s", err)
					}
					open = false
				}
			}
		}
	}()
}
