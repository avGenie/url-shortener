package entity

import (
	"fmt"
	"net/http"
)

const (
	UserIDKey = "user_id"
)

type UserIDCtxKey struct{}

type UserID string

func (u UserID) String() string {
	return string(u)
}

func (u UserID) IsValid() bool {
	return len(u.String()) != 0
}

func ValidateCookieUserID(cookie *http.Cookie) (string, error) {
	rawUserID := cookie.Value

	if len(rawUserID) == 0 {
		return "", fmt.Errorf("cookie of user id is empty")
	}

	return rawUserID, nil
}
