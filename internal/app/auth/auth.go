package authentication

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/google/uuid"
)

const (
	secretKey = "5269889d400bbf2dc66216f37b2839bb"
	idLength  = 128
	timeout   = 3 * time.Second
)

type UserAuthorisator interface {
	AddUser(ctx context.Context, userID entity.UserID) error
	AuthUser(ctx context.Context, userID entity.UserID) error
}

type UserAdder interface {
	AddUser(ctx context.Context, userID entity.UserID) error
}

type UserAuthenticator interface {
	AuthUser(ctx context.Context, userID entity.UserID) error
}

func GenerateCustomID() (string, error) {
	data := make([]byte, idLength)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(data)

	key := hash.Sum(data)

	return hex.EncodeToString(key), nil
}

func EncodeUserID(userID entity.UserID) (string, error) {
	rawData := userID.String()

	aes, err := aes.NewCipher([]byte(secretKey))
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

	aes, err := aes.NewCipher([]byte(secretKey))
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

func createUserID(userAdder UserAdder) (entity.UserID, error) {
	uuid := uuid.New()
	userID := entity.UserID(uuid.String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := userAdder.AddUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("error while add new user id to storage")
	}

	return userID, nil
}

func authenticateUser(userID entity.UserID, auth UserAuthenticator) error {
	ctx, close := context.WithTimeout(context.Background(), timeout)
	defer close()

	return auth.AuthUser(ctx, userID)
}
