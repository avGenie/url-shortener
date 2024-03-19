package get

import (
	"context"
	"net/http"
	"time"

	"github.com/avGenie/url-shortener/internal/app/entity"
	"github.com/avGenie/url-shortener/internal/app/handlers/errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type URLGetter interface {
	GetURL(ctx context.Context, key entity.URL) (*entity.URL, error)
}

// Processes GET request. Sends the source address at the given short address
//
// # Sends short URL back to the original using from the URL's map
//
// Returns 307 status code if processing was successfull, otherwise returns 400.
func URLHandler(getter URLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		shortURL := chi.URLParam(req, "url")

		ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
		defer cancel()

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

		url, err := getter.GetURL(ctx, *eShortURL)
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
