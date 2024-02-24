package handlers

import (
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/logger"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(config config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", logger.RequestLogger(PostMiddleware(config, PostHandlerURL)))
	r.Post("/api/shorten", logger.RequestLogger(PostMiddleware(config, PostHandlerJSON)))

	r.Get("/{url}", logger.RequestLogger(GetHandler))

	return r
}
