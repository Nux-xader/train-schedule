package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var SecretKey string

func Decrypt(data string, iv [16]byte) (result string, err error) {
	block, err := aes.NewCipher([]byte(SecretKey))
	if err != nil {
		return
	}

	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	cbc := cipher.NewCBCDecrypter(block, iv[:aes.BlockSize])

	plainText := make([]byte, len(cipherText))
	cbc.CryptBlocks(plainText, cipherText)

	result = string(plainText)

	length := len(result)
	if length == 0 {
		return
	}

	padLength := int(result[length-1])
	if padLength > length {
		return
	}

	result = result[:length-padLength]
	return
}
