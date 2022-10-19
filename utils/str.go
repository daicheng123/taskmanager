package utils

import "strings"

func BuilderStr(values ...string) string {
	var temp strings.Builder
	for _, value := range values {
		temp.WriteString(value)
	}
	return temp.String()
}
