package utils

import (
	"taskmanager/internal/consts"
	"taskmanager/internal/mapper"
	"taskmanager/pkg/logger"
	"time"
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
	if session.ExpireTime.After(a) && dur < consts.SessionCookieAge {
		// token 续期
		if dur < 300 {
			go RunSafeWithMsg(func() {
				session.ExpireTime = session.ExpireTime.Add(consts.SessionCookieAge * time.Second)
				sm.Save(session)
			}, "token 续期失败")
		}
		return true
	}
	return false
}
