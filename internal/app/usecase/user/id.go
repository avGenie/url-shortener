package user

import (
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/google/uuid"
)

// CreateUserID Creates user ID in UUID format
func CreateUserID() entity.UserID {
	id := uuid.New()
	userID := entity.UserID(id.String())

	return userID
}
