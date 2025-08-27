package http

import "github.com/go-chi/chi/v5"

func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/order/{orderUID}", handler.GetOrder)

	return r
}
