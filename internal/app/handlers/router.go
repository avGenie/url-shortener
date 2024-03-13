package handlers

import (
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/encoding"
	"github.com/avGenie/url-shortener/internal/app/handlers/get"
	"github.com/avGenie/url-shortener/internal/app/handlers/post"
	"github.com/avGenie/url-shortener/internal/app/logger"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api/model"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(config config.Config, db storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.LoggerMiddleware)
	r.Use(encoding.GzipMiddleware)

	r.Post("/", post.URLHandler(db, config.BaseURIPrefix))
	r.Post("/api/shorten", post.JSONHandler(db, config.BaseURIPrefix))
	r.Post("/api/shorten/batch", post.JSONBatchHandler(db, config.BaseURIPrefix))

	r.Get("/{url}", get.URLHandler(db))
	r.Get("/ping", get.PingDBHandler(db))

	return r
}
