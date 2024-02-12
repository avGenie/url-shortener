package handlers

import (
	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/go-chi/chi/v5"
)

func CreateRouter(config config.Config) *chi.Mux {
	r := chi.NewRouter()

	postContext := PostContext{
		baseURIPrefix: config.BaseURIPrefix,
		handle:        PostHandler,
	}

	r.Post("/", postContext.Handle())
	r.Get("/{url}", GetHandler)

	return r
}
