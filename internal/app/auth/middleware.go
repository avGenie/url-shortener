package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/entity"
	"go.uber.org/zap"
)

var (
	ErrInvalidRawUserID = errors.New("invalid raw user id")
)

func AuthMiddleware(config config.Config, userAuth UserAuthorisator) func(next http.Handler) http.Handler {
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

			userID, err := DecodeUserID(rawUserID)
			if err != nil {
				zap.L().Error("error while decoding user id in user authentication", zap.Error(err))
				if !errors.Is(err, ErrInvalidRawUserID) {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				processInvalidUserID(w, r, userAuth)
				return
			}

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

	encodedUserID, err := EncodeUserID(userID)
	if err != nil {
		zap.L().Error("error while encoding user id in user authentication", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.AddCookie(&http.Cookie{
		Name:  entity.UserIDKey,
		Value: encodedUserID,
	})
	w.WriteHeader(http.StatusUnauthorized)
}
