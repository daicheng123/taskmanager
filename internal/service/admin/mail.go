package admin

import (
	"fmt"
	"taskmanager/internal/cache"
	cacheutils "taskmanager/internal/cache/utils"
	"taskmanager/internal/consts"
	"taskmanager/internal/mapper"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
	"taskmanager/utils/serializer"
	"time"
)

type MailService struct {
	Email string `json:"email" binding:"required,email"`
}

func (ms *MailService) CheckMailExists(mail string) *serializer.Response {
	user, err := mapper.GetUserMapper().FindByEmail(mail)
	if err != nil {
		return serializer.DBErr("查询用户信息失败", err)
	}
	return &serializer.Response{Data: user}
}

func (ms *MailService) GenMailCode() *serializer.Response {
	so := cache.NewStringOperation()
	flagKey := ms.Email + "_" + consts.EmailFlagPrefix

	flagValue := so.Exists(flagKey).UnwrapOrElse(func(err error) {
		logger.Error("获取 key: %s 是否存在失败, err:[%s]", flagKey, err)
	})

	if flagValue.(int64) == 1 {
		return serializer.Err(serializer.CodeMailSendErr, "获取验证码操作频繁", nil)
	}

	code := utils.RandStringBytesMask(10)

	if r := so.Set(ms.Email, code, cacheutils.WithExpire(300*time.Second)).
		UnwrapOrElse(func(err error) {
			logger.Error("set key: %s 失败,err:[%s]", ms.Email, err)
		}); r == nil {
		return serializer.Err(serializer.CodeMailSendErr, "设置验证码失败", nil)
	}

	if r := so.Set(flagKey, true, cacheutils.WithExpire(60*time.Second)).
		UnwrapOrElse(func(err error) {
			logger.Error("set key: %s 失败,err:[%s]", flagKey, err)
		}); r == nil {
		return serializer.Err(serializer.CodeMailSendErr, "设置验证码失败", nil)
	}

	err := new(utils.MailService).
		Builder([]string{ms.Email}, consts.MailSubject).
		SetMsgBody(consts.HtmlMail, fmt.Sprintf(consts.MailCodeContent, code)).
		Build().
		Sender()

	if err != nil {
		logger.Error("验证码邮件发送失败,error:[%s]", err.Error())
		return serializer.Err(serializer.CodeMailSendErr, "验证码邮件发送失败", err)
	}
	return &serializer.Response{Message: "验证码推送成功！"}
}
