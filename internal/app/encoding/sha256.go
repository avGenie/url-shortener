package encoding

import "crypto/sha256"

// NewSHA256 Returns encoded byte array by given using SHA256 algorithm
func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)

	return hash[:]
}
