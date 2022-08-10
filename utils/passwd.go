package utils

import "golang.org/x/crypto/bcrypt"

func EncodePassword(pwdStr string) (pwdHash string, err error) {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	pwdHash = string(hash)
	return
}

func DecodeHashPassword(hashedPwd string, plainPwd string) bool {
	hashBytes := []byte(hashedPwd)
	strByte := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(hashBytes, strByte)
	if err != nil {
		return false
	}
	return true
}
