package random

import (
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomURL() string {
	return fmt.Sprintf("https://%s.com", generateRandomString(5, 15))
}

func GenerateURL(randomString string) string {
	return fmt.Sprintf("https://%s.com", randomString)
}

func GenerateRandomString() string {
	return generateRandomString(5, 15)
}

func SleepDuration(min, max int) time.Duration {
	n := rand.Intn(max-min) + min
	return time.Duration(n)
}

func generateRandomString(minLen, maxLen int) string {
	n := rand.Intn(maxLen-minLen) + minLen
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
