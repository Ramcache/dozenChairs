package main

import (
	"dozenChairs/pkg/app"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/db"
	"log"
	"net/http"

	_ "dozenChairs/docs"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// @title DozenChairs API
// @version 1.0
// @description API для магазина мебели
// @BasePath /
// @schemes http
func main() {
	cfg := config.Load()

	_ = godotenv.Load()
	pool, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}
	defer pool.Close() // defet
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	r := app.InitRoutes(pool)

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", corsMiddleware.Handler(r))
}
