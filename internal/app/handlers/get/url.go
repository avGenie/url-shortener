package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/handlers/errors"
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

		userID, ok := req.Context().Value(entity.UserIDCtxKey{}).(entity.UserID)
		if !ok {
			zap.L().Error("user id couldn't obtain from context")
			writer.WriteHeader(http.StatusInternalServerError)
			return
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
			zap.L().Error(
				"error while getting url",
				zap.String("error", err.Error()),
				zap.String("short_url", shortURL),
			)

			http.Error(writer, errors.ShortURLNotInDB, http.StatusBadRequest)
			return
		}

		zap.L().Info("url has been decoded succeessfully", zap.String("decoded url", url.String()))

		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.Header().Set("Location", url.String())
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func UserURLsHandler(getter AllURLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(entity.UserIDCtxKey{}).(entity.UserID)
		if !ok {
			zap.L().Error("user id couldn't obtain from context while all user urls processing")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		urls, err := getter.GetAllURLByUserID(ctx, userID)
		if err != nil {
			zap.L().Error("couldn't get all user urls", zap.Error(err), zap.String("user_id", userID.String()))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		out, err := json.Marshal(urls)
		if err != nil {
			zap.L().Error("error while converting all user urls to output", zap.Error(err))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(out)
	}
}
