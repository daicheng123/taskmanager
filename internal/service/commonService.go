package service

import (
	"fmt"
	"taskmanager/internal/cache"
	cacheutils "taskmanager/internal/cache/utils"
	"taskmanager/internal/consts"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
	"time"
)

type CommonService struct {
}

func NewCommonService() *CommonService {
	return &CommonService{}
}

func (us *CommonService) GenEmailCode(email string) error {
	if email == "" {
		return fmt.Errorf("邮件地址不能为空！")
	}
	so := cache.NewStringOperation()
	flagKey := email + "_" + consts.EmailFlagPrefix

	flagValue := so.Exists(flagKey).UnwrapOrElse(func(err error) {
		logger.Error("获取 key: %s 是否存在失败, err:[%s]", flagKey, err)
	})

	if flagValue.(int64) == 1 {
		return fmt.Errorf("获取验证码操作频繁！")
	}

	code := utils.RandStringBytesMask(10)

	if r := so.Set(email, code, cacheutils.WithExpire(300*time.Second)).
		UnwrapOrElse(func(err error) {
			logger.Error("set key: %s 失败,err:[%s]", email, err)
		}); r == nil {
		return fmt.Errorf("设置验证码失败")
	}

	if r := so.Set(flagKey, true, cacheutils.WithExpire(60*time.Second)).
		UnwrapOrElse(func(err error) {
			logger.Error("set key: %s 失败,err:[%s]", flagKey, err)
		}); r == nil {
		return fmt.Errorf("设置验证码失败")
	}

	err := new(utils.MailService).
		Builder([]string{email}, consts.MailSubject).
		SetMsgBody(consts.HtmlMail, fmt.Sprintf(consts.MailCodeContent, code)).
		Build().
		Sender()

	if err != nil {
		logger.Error("验证码邮件发送失败,error:[%s]", err.Error())
		return err
	}
	return nil
}
