package service

import (
	"taskmanager/internal/dal/mapper"
	"taskmanager/pkg/logger"
	"taskmanager/utils"
	"time"
)

const (
	SessionCookieAge = 7200
)

//SessionJudge  校验token
func SessionJudge(token string) bool {
	sm := mapper.GetSessionMapper()
	session, err := sm.FindByToken(token)
	if err != nil {
		logger.GinERROR("获取会话失败, err:[%s]", err.Error())
		return false
	}
	dur := session.ExpireTime.Sub(time.Now()).Seconds()
	a := time.Now()
	if session.ExpireTime.After(a) && dur < SessionCookieAge {
		// token 续期
		if dur < 300 {
			go utils.RunSafeWithMsg(func() {
				session.ExpireTime = session.ExpireTime.Add(SessionCookieAge * time.Second)
				sm.Save(session)
			}, "token 续期失败")
		}
		return true
	}
	return false
}
