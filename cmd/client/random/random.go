// Package random generates random URL
package random

import (
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomURL Generates random URL
func GenerateRandomURL() string {
	return fmt.Sprintf("https://%s.com", generateRandomString(5, 15))
}

// GenerateURL Generates URL using random input string
func GenerateURL(randomString string) string {
	return fmt.Sprintf("https://%s.com", randomString)
}

// GenerateRandomString Generates random domain
func GenerateRandomString() string {
	return generateRandomString(5, 15)
}

// SleepDuration Generates random time duration between given min and max
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
