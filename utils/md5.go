package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

func GenMd5(value string) string {
	salt := strconv.FormatInt(time.Now().UnixNano(), 10)
	d := []byte(value)
	m := md5.New()
	m.Write([]byte(salt))
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
