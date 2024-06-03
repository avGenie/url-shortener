package entity

import (
	"fmt"
	"net/http"
)

// UserIDKey Cookie key to store user ID
const UserIDKey = "user_id"

// UserIDCtxKey Key to store user ID in go context
type UserIDCtxKey struct{}

// UserIDCtx Value to store user ID in go context
type UserIDCtx struct {
	UserID     UserID
	StatusCode int
}

// UserID Contains user ID
type UserID string

// String Implements stringer interface
func (u UserID) String() string {
	return string(u)
}

// IsValid Validates user ID
func (u UserID) IsValid() bool {
	return len(u.String()) != 0
}

// ValidateCookieUserID Validates user ID obtained from cookie
func ValidateCookieUserID(cookie *http.Cookie) (UserID, error) {
	rawUserID := cookie.Value

	if len(rawUserID) == 0 {
		return "", fmt.Errorf("cookie of user id is empty")
	}

	return UserID(rawUserID), nil
}
