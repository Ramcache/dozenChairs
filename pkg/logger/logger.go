package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(isProduction bool) {
	// создаём директорию logs, если не существует
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatalf("failed to create logs directory: %v", err)
		}
	}

	// Настраиваем лог-файл
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.Lock(os.Stdout)

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, zap.DebugLevel),
		zapcore.NewCore(encoder, consoleWriter, zap.DebugLevel),
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
