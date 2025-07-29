package app

import (
	_ "dozenChairs/docs"
	"dozenChairs/internal/auth"
	"dozenChairs/internal/handlers"
	"dozenChairs/internal/middlewares"
	"dozenChairs/internal/repository"
	"dozenChairs/internal/services"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	r chi.Router,
	productHandler *handlers.ProductHandler,
	authHandler *handlers.AuthHandler,
	imageHandler *handlers.ImageHandler,
	jwtManager *auth.JWTManager,
) {
	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {

		// --- Public ---
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", authHandler.Register)
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/refresh", authHandler.Refresh)
			r.Post("/auth/logout", authHandler.Logout)

			r.Get("/products", productHandler.GetAll)
			r.Get("/products/{slug}", productHandler.GetBySlug)
			r.Get("/products/sets/{slug}", productHandler.GetSetBySlug)

			r.Get("/sets", productHandler.GetSets)
			r.Get("/categories", productHandler.GetCategories)

			// Публичный просмотр изображений по товару
			r.Get("/products/{product_id}/images", imageHandler.GetByProductID)
		})

		// --- Authorized Users ---
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireAuth(jwtManager))
			r.Get("/auth/me", authHandler.Me)
		})

		// --- Admin-only ---
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireAuth(jwtManager))
			r.Use(middlewares.RequireRole("admin"))

			// Товары
			r.Post("/products", productHandler.Create)
			r.Put("/products/{slug}", productHandler.Update)
			r.Delete("/products/{slug}", productHandler.Delete)

			// Изображения
			r.Post("/upload", imageHandler.Upload)
			r.Delete("/images/{id}", imageHandler.Delete)
		})
	})
}

func SetupRouter(cfg *config.Config, log logger.Logger, conn *pgxpool.Pool) http.Handler {
	// Репозитории
	userRepo := repository.NewUserRepo(conn)
	sessionRepo := repository.NewSessionRepo(conn)
	imageRepo := repository.NewImageRepo(conn)
	productRepo := repository.NewProductRepo(conn)

	// Сервисы
	authService := services.NewAuthService(userRepo, sessionRepo)
	imageService := services.NewImageService(imageRepo)
	productService := services.NewProductService(productRepo)

	// JWT
	jwtManager := auth.NewJWTManager(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)

	// Хендлеры
	authHandler := handlers.NewAuthHandler(authService, log, jwtManager)
	imageHandler := handlers.NewImageHandler(imageService)
	productHandler := handlers.NewProductHandler(productService, log)

	// Роутер
	r := chi.NewRouter()
	r.Use(middlewares.Recover(log))
	r.Use(middlewares.RequestID())
	r.Use(middlewares.CORS())
	r.Use(middlewares.Logger(log))

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	RegisterRoutes(r, productHandler, authHandler, imageHandler, jwtManager)

	return r
}
