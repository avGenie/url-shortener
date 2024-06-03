// Package auth provides authenticate middleware and encodes and decodes user ID
package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/usecase/user"
	"go.uber.org/zap"
)

// AuthMiddleware authenticate middleware validates user ID obtained from cookies
//
// Returns 200(StatusOK) if validation was performed correctly
// Returns 401(StatusUnauthorized) if cookie with user ID undefined
// Returns 401(StatusUnauthorized) if user ID obtained from cookies is invalid
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("start user authentication")

		status := http.StatusOK
		userIDCookie, err := r.Cookie(entity.UserIDKey)

		// Cookies doesn't contain user id
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				zap.L().Info("cookie with user id is not defined")
			} else {
				zap.L().Info("error while getting cookie", zap.Error(err))
			}

			status = http.StatusUnauthorized

			userIDCookie = processInvalidCookie(w)
		}

		// User id invalid: may be empty
		userID, err := entity.ValidateCookieUserID(userIDCookie)
		if err != nil {
			zap.L().Error("error while validating user id from cookie in user authentication", zap.Error(err))
			status = http.StatusUnauthorized

			processInvalidCookie(w)
		}

		userCtx := entity.UserIDCtx{
			UserID:     userID,
			StatusCode: status,
		}

		ctx := context.WithValue(r.Context(), entity.UserIDCtxKey{}, userCtx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func processInvalidCookie(w http.ResponseWriter) *http.Cookie {
	userID := user.CreateUserID()

	cookie := &http.Cookie{
		Name:  entity.UserIDKey,
		Value: userID.String(),
	}

	http.SetCookie(w, cookie)

	return cookie
}
