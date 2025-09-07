package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *Handler, staticFS http.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/order/{orderUID}", handler.GetOrder)

	r.Handle("/static/*", http.StripPrefix("/static/", staticFS))

	return r
}
