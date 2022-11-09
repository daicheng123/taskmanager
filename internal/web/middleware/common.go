package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

type MiddleWare interface {
	OnRequest() gin.HandlerFunc
}

var NoAuthUri = []string{
	"/health",
	"/check_email_exists",
	"/email_code",
	"/users/login",
	"/users/register",
	"/users/user_code",
	"/websockets/core",
}

func AllowUri(uri string) bool {
	var reg = &regexp.Regexp{}
	if strings.HasPrefix(uri, "/api/v1") {
		reg = regexp.MustCompile("^/api/v1(.*)$")
	} else if strings.HasPrefix(uri, "/api/") {
		reg = regexp.MustCompile("^/api(.*)$")
	} else {
		reg = regexp.MustCompile("(.*)")
	}
	uriStr := reg.ReplaceAllString(uri, "$1")
	for _, u := range NoAuthUri {
		if strings.Contains(uriStr, u) {
			return true
		}
	}
	return false
}
