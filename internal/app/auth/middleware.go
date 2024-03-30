package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrInvalidRawUserID = errors.New("invalid raw user id")
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

func AuthMiddleware(userAuth UserAuthorisator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			zap.L().Info("start user authentication")

			userIDCookie, err := r.Cookie(entity.UserIDKey)

			// Cookies doesn't contain user id
			if err != nil {
				zap.L().Debug("error while getting cookie with user id in user authentication", zap.Error(err))
				processInvalidUserID(w, r, userAuth)
				return
			}

			// User id invalid: may be empty
			rawUserID, err := entity.ValidateCookieUserID(userIDCookie)
			if err != nil {
				zap.L().Error("error while validating user id from cookie in user authentication", zap.Error(err))
				processInvalidUserID(w, r, userAuth)
				return
			}

			userID := entity.UserID(rawUserID)

			// userID, err := DecodeUserID(rawUserID)
			// if err != nil {
			// 	zap.L().Error("error while decoding user id in user authentication", zap.Error(err))
			// 	if !errors.Is(err, ErrInvalidRawUserID) {
			// 		w.WriteHeader(http.StatusInternalServerError)
			// 		return
			// 	}

			// 	processInvalidUserID(w, r, userAuth)
			// 	return
			// }

			// User id is not in DB
			err = authenticateUser(userID, userAuth)
			if err != nil {
				zap.L().Error("error while user authentication", zap.Error(err))
				// TODO: add condition if error is in db
				processInvalidUserID(w, r, userAuth)
				return
			}

			ctx := context.WithValue(context.Background(), entity.UserIDCtxKey{}, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func processInvalidUserID(w http.ResponseWriter, r *http.Request, userAuth UserAuthorisator) {
	userID, err := createUserID(userAuth)

	if err != nil {
		zap.L().Error("error while creating user id in user authentication", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encodedUserID, err := EncodeUserID(userID)
	// if err != nil {
	// 	zap.L().Error("error while encoding user id in user authentication", zap.Error(err))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	r.AddCookie(&http.Cookie{
		Name:   entity.UserIDKey,
		Value:  userID.String(),
		Secure: true,
	})
	w.WriteHeader(http.StatusUnauthorized)
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
