package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	handler_err "github.com/avGenie/url-shortener/internal/app/handlers/errors"
	storage_err "github.com/avGenie/url-shortener/internal/app/storage/api/errors"
	"github.com/avGenie/url-shortener/internal/app/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(ctx context.Context, userID entity.UserID, key entity.URL) (*entity.URL, error)
}

type AllURLGetter interface {
	GetAllURLByUserID(ctx context.Context, userID entity.UserID) (models.AllUrlsBatch, error)
}

// Processes GET request. Sends the source address at the given short address
//
// # Sends short URL back to the original using from the URL's map
//
// Returns 307 status code if processing was successfull, otherwise returns 400.
func URLHandler(getter URLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		shortURL := chi.URLParam(req, "url")

		var userID entity.UserID
		userIDCtx, ok := req.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)
		if ok {
			if userIDCtx.StatusCode == http.StatusOK {
				userID = userIDCtx.UserID
			} else {
				zap.L().Info("user id couldn't obtain from context")
			}
		} else {
			zap.L().Info("user id is empty from context")
		}

		eShortURL, err := entity.ParseURL(shortURL)
		if err != nil {
			zap.L().Error(
				"error while parsing short url",
				zap.String("error", err.Error()),
				zap.String("short_url", shortURL),
			)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), pingTimeout)
		defer cancel()

		url, err := getter.GetURL(ctx, userID, *eShortURL)
		if err != nil {
			if errors.Is(err, storage_err.ErrAllURLsDeleted) {
				writer.WriteHeader(http.StatusGone)
				return
			}

			zap.L().Error(
				"error while getting url",
				zap.String("error", err.Error()),
				zap.String("short_url", shortURL),
			)

			http.Error(writer, handler_err.ShortURLNotInDB, http.StatusBadRequest)
			return
		}

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.Header().Set("Location", url.String())
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func UserURLsHandler(getter AllURLGetter, baseURIPrefix string) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		if baseURIPrefix == "" {
			zap.L().Error("invalid base URI prefix", zap.String("base URI prefix", baseURIPrefix))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		userIDCtx, ok := req.Context().Value(entity.UserIDCtxKey{}).(entity.UserIDCtx)
		if !ok {
			zap.L().Error("user id couldn't obtain from context while all user urls processing")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if code := validateUserIDCtx(userIDCtx); code != http.StatusOK {
			writer.WriteHeader(code)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		out, err := processAllUSerURL(getter, ctx, userIDCtx.UserID, baseURIPrefix)
		if err != nil {
			if errors.Is(err, ErrAllURLNotFound) {
				writer.WriteHeader(http.StatusNoContent)
				return
			}

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(out)
	}
}

func validateUserIDCtx(userIDCtx entity.UserIDCtx) int {
	if userIDCtx.StatusCode == http.StatusUnauthorized {
		zap.L().Error("user id couldn't obtain from context")
		return userIDCtx.StatusCode
	}

	if len(userIDCtx.UserID.String()) == 0 {
		zap.L().Error("empty user id from context")
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
