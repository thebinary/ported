package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

const randomServiceNameLen = 15

func generateServiceName() (serviceName string) {
	currentTimeString := strconv.Itoa(int(time.Now().UnixNano()))
	sha := sha256.Sum256([]byte(currentTimeString))
	serviceName = fmt.Sprintf("%x", sha)[:randomServiceNameLen]
	return serviceName
}
