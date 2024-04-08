package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
)

const (
	secretKey = "5269889d400bbf2dc66216f37b2839bb"
	idLength  = 128
	timeout   = 3 * time.Second
)

func EncodeUserID(userID entity.UserID) (string, error) {
	rawData := userID.String()

	decodedKey, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error while decoding user id in creating cipher process: %w", err)
	}

	aes, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", fmt.Errorf("error while encoding user id in creating cipher process: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("error while encoding user id in creating GCM process: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("error while encoding user id: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(rawData), nil)

	return string(ciphertext), nil
}

func DecodeUserID(data string) (entity.UserID, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("error while decoding user id: %w", ErrInvalidRawUserID)
	}

	decodedKey, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error while decoding user id in creating cipher process: %w", err)
	}

	aes, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", fmt.Errorf("error while decoding user id in creating cipher process: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("error while decoding user id in creating GCM process: %w", err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	userID, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", fmt.Errorf("error while decoding user id: %w", err)
	}

	return entity.UserID(string(userID)), nil
}
