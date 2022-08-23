package consts

const (
	UserControllerGroup = "user"
	AppManagerConfPath  = "APP_MANAGER_CONF_FILE"

	ManagerLog = iota // 服务日志
	GinLog            // gin框架日志
	TaskLog           // 任务日志

	EmailFlagPrefix = "send"

	AttrExpire = "expire"

	HtmlMail uint8 = 1
	TextMail uint8 = 2

	SessionCookieAge = 7200

	UserTokenStr = "USER_TOKEN"

	MailSubject     = "【运维任务管理平台】"
	MailCodeContent = "您的注册验证为：<span style='color:red'>%s</span>，请在5分钟内注册。"
)
