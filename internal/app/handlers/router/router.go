package handlers

import (
	"github.com/avGenie/url-shortener/internal/app/auth"
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/encoding"
	handlers "github.com/avGenie/url-shortener/internal/app/handlers/delete"
	get "github.com/avGenie/url-shortener/internal/app/handlers/get"
	post "github.com/avGenie/url-shortener/internal/app/handlers/post"
	"github.com/avGenie/url-shortener/internal/app/logger"
	storage "github.com/avGenie/url-shortener/internal/app/storage/api/model"
	cidr "github.com/avGenie/url-shortener/internal/app/usecase/CIDR"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router Routes endpoint handlers
type Router struct {
	Mux *chi.Mux

	deleteHandler *handlers.DeleteHandler
}

// NewRouter Creates router
func NewRouter(config config.Config, db storage.Storage, cidr *cidr.CIDR) *Router {
	deleteHandler := handlers.NewDeleteHandler(db)
	return &Router{
		Mux:           createRouter(config, deleteHandler, db, cidr),
		deleteHandler: deleteHandler,
	}
}

// Stop Stops router
func (r *Router) Stop() {
	r.deleteHandler.Stop()
}

func createRouter(
	config config.Config,
	deleteHandler *handlers.DeleteHandler,
	db storage.Storage,
	cidr *cidr.CIDR,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.LoggerMiddleware)
	r.Use(encoding.GzipMiddleware)
	r.Use(auth.AuthMiddleware)

	r.Mount("/debug", middleware.Profiler())

	r.Post("/", post.URLHandler(db, config.BaseURIPrefix))
	r.Post("/api/shorten", post.JSONHandler(db, config.BaseURIPrefix))
	r.Post("/api/shorten/batch", post.JSONBatchHandler(db, config.BaseURIPrefix))

	r.Get("/{url}", get.URLHandler(db))
	r.Get("/ping", get.PingDBHandler(db))
	r.Get("/api/internal/stats", get.StatsHandler(db, cidr))
	r.Get("/api/user/urls", get.UserURLsHandler(db, config.BaseURIPrefix))

	r.Delete("/api/user/urls", deleteHandler.DeleteUserURLHandler())

	return r
}
