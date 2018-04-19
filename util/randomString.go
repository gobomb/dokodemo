package util

import (
	"github.com/qiniu/log"
	"github.com/satori/go.uuid"
)

func GenerateRandomString() string {
	randomString, err := uuid.NewV4()
	if err != nil {
		log.Errorf("[generateRandomString]create fail %v", err)
	}
	return randomString.String()
}
