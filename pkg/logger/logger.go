package logger

import (
	"log"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(isProduction bool) {
	var err error

	if isProduction {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("failed to init zap logger: %v", err)
	}
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
