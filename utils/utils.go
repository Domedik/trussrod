package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
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

func LimitAndOffset(query url.Values) (uint, uint) {
	queryPage := query["page"]
	var p string
	if len(queryPage) > 0 {
		p = queryPage[0]
	}
	querySize := query["size"]
	var size string
	if len(querySize) > 0 {
		size = querySize[0]
	}
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
