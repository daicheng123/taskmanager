package utils

import (
	"bufio"
	"io/ioutil"
	"os"
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

func WriteFile(filePath string, contents string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// 创建一个写的对象
	wr := bufio.NewWriter(file)
	defer wr.Flush()
	_, err = wr.WriteString(contents)
	return err
}
