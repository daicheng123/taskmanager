package admin

import (
	"context"
	"fmt"
	"taskmanager/internal/cache"
	cacheutils "taskmanager/internal/cache/utils"
	"taskmanager/internal/repo/mapper"
	"taskmanager/pkg/email"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/serializer"
	"taskmanager/utils"
	"time"
)

const (
	EmailFlagPrefix = "send"
	MailCodeContent = "您的注册验证为：<span style='color:red'>%s</span>，请在5分钟内注册。"
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

func (ms *MailService) GenMailCode(ctx context.Context) *serializer.Response {
	rc := cache.NewRedisCache(ctx)
	flagKey := utils.BuilderStr(ms.Email, "_", EmailFlagPrefix)
	//flagValue := rc.Exists(flagKey).UnwrapOrElse(func(err error) {
	//	logger.Error("获取 key: %s 是否存在失败, err:[%s]", flagKey, err)
	//})
	flagValue, err := rc.Exists(flagKey)
	if err != nil {
		logger.Error("获取 key: %s 是否存在失败, err:[%s]", flagKey, err)
		return serializer.Err(serializer.CodeMailSendErr, "获取邮件验证码发送标识位失败", err)
		//return serializer.Err()
	}

	if flagValue == 1 {
		return serializer.Err(serializer.CodeMailSendErr, "获取验证码操作频繁", nil)
	}

	code := utils.RandStringBytesMask(10)

	if err = rc.Set(ms.Email, code, cacheutils.WithExpire(300*time.Second)); err != nil {
		//logger.Error("set key: %s 失败,err:[%s]", ms.Email, err)
		return serializer.Err(serializer.CodeMailSendErr, "设置验证码失败", err)
	}

	if err = rc.Set(flagKey, true, cacheutils.WithExpire(60*time.Second)); err != nil {
		return serializer.Err(serializer.CodeMailSendErr, "设置验证码失败", nil)
	}

	//if r := rc.Set(ms.Email, code, cacheutils.WithExpire(300*time.Second)).
	//	UnwrapOrElse(func(err error) {
	//		logger.Error("set key: %s 失败,err:[%s]", ms.Email, err)
	//	}); r == nil {
	//	return serializer.Err(serializer.CodeMailSendErr, "设置验证码失败", nil)
	//}

	//if r := so.Set(flagKey, true, cacheutils.WithExpire(60*time.Second)).
	//	UnwrapOrElse(func(err error) {
	//		logger.Error("set key: %s 失败,err:[%s]", flagKey, err)
	//	}); r == nil {
	//	return serializer.Err(serializer.CodeMailSendErr, "设置验证码失败", nil)
	//}

	//err := new(email.MailService).
	//	Builder([]string{ms.Email}, email.MailSubject).
	//	SetMsgBody(email.HtmlMail, fmt.Sprintf(MailCodeContent, code)).
	//	Build().
	//	Sender()
	err = email.Send([]string{ms.Email}, "", fmt.Sprintf(MailCodeContent, code))
	if err != nil {
		logger.Error("验证码邮件发送失败,error:[%s]", err.Error())
		return serializer.Err(serializer.CodeMailSendErr, "验证码邮件发送失败", err)
	}
	return &serializer.Response{Message: "验证码推送成功！"}
}
