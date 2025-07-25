package app

import (
	"dozenChairs/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *handlers.ProductHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Post("/", h.Create)
			r.Get("/", h.GetAll)
			r.Get("/{slug}", h.GetBySlug)
			r.Get("/sets/{slug}", h.GetSetBySlug)
			r.Put("/{slug}", h.Update)
			r.Delete("/{slug}", h.Delete)
		})

		r.Get("/sets", h.GetSets)
		r.Get("/categories", h.GetCategories)
	})
}
