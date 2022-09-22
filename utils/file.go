package utils

import (
	"io/ioutil"
	"strings"
)

func ReadFile(filePath string) (string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fileStr := string(file)
	return strings.Trim(fileStr, "\r\t\n"), err
}
