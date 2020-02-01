package main

import (
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
)

func Hash(originalUrl string) []byte {
	h := sha256.New()
	h.Write([]byte(originalUrl))
	hasResult := h.Sum(nil)
	return hasResult
}

func Encode(hashResult []byte) string {
	//TODO: keep generating if conflict now, change it later
	encodedStr := base58.Encode(hashResult)[:8]
	return encodedStr
}

func GenerateKey(originalUrl string) string {
	hashRes := Hash(originalUrl)
	encodedUrl := Encode(hashRes)
	return encodedUrl
}
