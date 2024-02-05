package handlers

import "github.com/go-chi/chi/v5"

func CreateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", PostHandler)
	r.Get("/{url}", GetHandler)

	return r
}
