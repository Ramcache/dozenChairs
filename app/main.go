package main

import (
	"dozenChairs/pkg/app"
	"dozenChairs/pkg/config"
	"dozenChairs/pkg/db"
	"dozenChairs/pkg/httphelper"
	"dozenChairs/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	// Инициализация логгера
	logger.Init(false)
	defer logger.Sync()
	log := logger.Log

	// Подключение к БД
	conn := db.MustConnectDB(cfg, log)
	defer conn.Close()

	// Сборка зависимостей и роутера
	r := app.SetupRouter(cfg, log, conn)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Старт сервера в горутине
	go func() {
		log.Info("server started", zap.String("port", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	httphelper.WaitForShutdown(srv, log)
}
