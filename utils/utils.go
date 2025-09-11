package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
)

func GetSecretHash(username, secret string) string {
	data := username
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	hash := h.Sum(nil)
	secretHash := base64.StdEncoding.EncodeToString(hash)
	return secretHash
}

func LimitAndOffset(p, size string) (uint, uint) {
	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		page = 1
	}
	l, err := strconv.Atoi(size)
	if err != nil || l < 1 {
		l = 10
	}

	offset := (page - 1) * l
	return uint(l), uint(offset)
}
