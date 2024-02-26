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

	r.Post("/", PostMiddleware(config, PostHandlerURL))
	r.Post("/api/shorten", PostMiddleware(config, PostHandlerJSON))

	r.Get("/{url}", GetHandler)

	return r
}
