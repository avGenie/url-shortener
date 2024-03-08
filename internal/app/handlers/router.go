package handlers

import (
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/encoding"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(config config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.LoggerMiddleware)
	r.Use(encoding.GzipMiddleware)

	postContextURL := NewPostContextURL(config)
	postContextJSON := NewPostContextJSON(config)

	r.Post("/", postContextURL.Handle())
	r.Post("/api/shorten", postContextJSON.Handle())

	r.Get("/{url}", GetHandler)

	return r
}
