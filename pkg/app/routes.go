package app

import (
	"dozenChairs/internal/auth"
	"dozenChairs/internal/handlers"
	"dozenChairs/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r chi.Router,
	productHandler *handlers.ProductHandler,
	authHandler *handlers.AuthHandler,
	jwtManager *auth.JWTManager,
) {
	r.Route("/api/v1", func(r chi.Router) {
		// --- Public Auth ---
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Post("/auth/logout", authHandler.Logout)

		// --- Public Product Access ---
		r.Get("/products", productHandler.GetAll)
		r.Get("/products/{slug}", productHandler.GetBySlug)
		r.Get("/products/sets/{slug}", productHandler.GetSetBySlug)

		r.Get("/sets", productHandler.GetSets)
		r.Get("/categories", productHandler.GetCategories)

		// --- Authenticated Users ---
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireAuth(jwtManager))
			r.Get("/auth/me", authHandler.Me)
		})

		// --- Admin Only ---
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireAuth(jwtManager))
			r.Use(middlewares.RequireRole("admin"))

			r.Post("/products", productHandler.Create)
			r.Put("/products/{slug}", productHandler.Update)
			r.Delete("/products/{slug}", productHandler.Delete)
		})
	})
}
