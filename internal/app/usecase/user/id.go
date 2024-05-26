package user

import (
	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/google/uuid"
)

func CreateUserID() entity.UserID {
	id := uuid.New()
	userID := entity.UserID(id.String())

	return userID
}
