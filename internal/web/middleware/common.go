package middleware

import (
	"github.com/gin-gonic/gin"
	"reflect"
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
	"/user/login",
	"/user/register",
	"/user/code",
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
		if reflect.DeepEqual(u, uriStr) {
			return true
		}
	}
	return false
}
