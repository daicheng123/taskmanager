package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"taskmanager/pkg/serializer"
)

//func EncodePassword(pwdStr string) (pwdHash string, err error) {
//	pwd := []byte(pwdStr)
//	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
//	if err != nil {
//		return
//	}
//	pwdHash = string(hash)
//	return
//}
//
//func CompareHashPassword(hashedPwd string, plainPwd string) bool {
//	hashBytes := []byte(hashedPwd)
//	strByte := []byte(plainPwd)
//	err := bcrypt.CompareHashAndPassword(hashBytes, strByte)
//	if err != nil {
//		return false
//	}
//	return true
//}

var (
	cipherKey               = []byte("0qi1lerodnvcm09181*^LkqDv,!dP>1U")
	ErrEncodePasswordFailed = serializer.NewError(serializer.CodeEncodePasswordErr, "密码编码失败", nil)
	ErrDecodePasswordFailed = serializer.NewError(serializer.CodeDecodePasswordErr, "密码解码失败", nil)
)

func Encrypt(message string) (encoded string, err error) {
	//Create byte array from the input string
	plainText := []byte(message)
	//Create a new AES cipher using the key
	block, err := aes.NewCipher(cipherKey)
	//IF NewCipher failed, exit:
	if err != nil {
		return "", ErrEncodePasswordFailed.WithError(err)
	}
	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", ErrEncodePasswordFailed.WithError(err)
	}
	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)
	//IF DecodeString failed, exit:
	if err != nil {
		return "", ErrDecodePasswordFailed.WithError(err)
	}
	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(cipherKey)
	//IF NewCipher failed, exit:
	if err != nil {
		return "", ErrDecodePasswordFailed.WithError(err)
	}
	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short!")
		return "", ErrDecodePasswordFailed.WithError(err)
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}
