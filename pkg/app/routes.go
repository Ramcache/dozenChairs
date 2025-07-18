package app

import (
	"dozenChairs/internal/handlers"
	middleware "dozenChairs/internal/middlewares"
	"dozenChairs/internal/repository"
	"dozenChairs/internal/services"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRoutes(pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	repo := repository.NewUserRepository(pool)
	service := services.NewUserService(repo)
	handler := handlers.NewAuthHandler(service)

	chairRepo := repository.NewChairRepository(pool)
	chairService := services.NewChairService(chairRepo)
	chairHandler := handlers.NewChairHandler(chairService)

	r.Post("/chairs", chairHandler.Create)
	r.Get("/chairs", chairHandler.GetAll)
	r.Get("/chairs/{slug}", chairHandler.GetBySlug)
	r.Patch("/chairs/{slug}", chairHandler.UpdateBySlug)
	r.Delete("/chairs/{slug}", chairHandler.DeleteBySlug)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)

	r.Group(func(priv chi.Router) {
		priv.Use(middleware.JWTAuth)
		priv.Get("/profile", handler.Profile)

		//priv.Group(func(admin chi.Router) {
		//	admin.Use(middleware.OnlyAdmin)
		//	admin.Post("/products", productHandler.Create)
		//	// ... другие admin-only ручки
		//})
	})
	return r
}
