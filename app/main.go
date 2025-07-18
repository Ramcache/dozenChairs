package main

import (
	"dozenChairs/pkg/app"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/db"
	"log"
	"net/http"

	_ "dozenChairs/docs"

	"github.com/joho/godotenv"
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

	r := app.InitRoutes(pool)

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", r)
}
