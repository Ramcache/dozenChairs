package main

import (
	"context"
	_ "dozenChairs/docs"
	"dozenChairs/internal/middlewares"
	"dozenChairs/pkg/app"
	"github.com/go-chi/chi/v5"
	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dozenChairs/internal/handlers"
	"dozenChairs/internal/repository"
	"dozenChairs/internal/services"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/db"
	"dozenChairs/pkg/logger"
)

// @title DozenChairs API
// @version 1.0
// @description REST API for managing chairs, tables and sets
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email ramaro@internet.ru

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.LoadConfig()
	logger.Init(false)
	defer logger.Sync()

	zapLog := logger.NewZapLogger(logger.Log)

	conn, err := db.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer conn.Close(context.Background())

	productRepo := repository.NewProductRepo(conn)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService, zapLog)

	r := chi.NewRouter()

	r.Use(middlewares.Recover(zapLog))
	r.Use(middlewares.RequestID())
	r.Use(middlewares.CORS())
	r.Use(middlewares.Logger(zapLog))
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	app.RegisterRoutes(r, productHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		logger.Log.Info("server started", zap.String("port", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("server shutdown failed", zap.Error(err))
	}

	logger.Log.Info("server exited properly")
}
