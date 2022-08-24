package common

import (
	"database/sql/driver"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"time"
)

// https://segmentfault.com/a/1190000022264001

const TimeFormat = "2006-01-02 15:04:05"

type CustomTime time.Time

func (t *CustomTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = CustomTime(time.Time{})
		return
	}

	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = CustomTime(now)
	return
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t CustomTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *CustomTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = CustomTime(tTime)
	return nil
}

func (t CustomTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

// ValidateJSONDateType 注册校验器
func ValidateJSONDateType(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(CustomTime{}) {
		timeStr := field.Interface().(CustomTime).String()
		// 0001-01-01 00:00:00 是 go 中 time.Time 类型的空值
		// 这里返回 Nil 则会被 validator 判定为空值，而无法通过 `binding:"required"` 规则
		if timeStr == "0001-01-01 00:00:00" {
			return nil
		}
		return timeStr
	}
	return nil
}

func TopicUrl(fl validator.FieldLevel) bool {
	if url, ok := fl.Field().Interface().(string); ok {

		if matched, _ := regexp.MatchString(`\w{4,10}`, url); matched {
			return true
		}
	}
	return false
}
