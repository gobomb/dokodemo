package util

import (
	"github.com/satori/go.uuid"
	"github.com/qiniu/log"
)

func GenerateRandomString() string {
	randomString, err := uuid.NewV4()
	if err != nil {
		log.Errorf("[generateRandomString]create fail %v", err)
	}
	return randomString.String()
}
