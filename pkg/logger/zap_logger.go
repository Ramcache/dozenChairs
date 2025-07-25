package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{l: l}
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field)  { z.l.Info(msg, fields...) }
func (z *ZapLogger) Error(msg string, fields ...zap.Field) { z.l.Error(msg, fields...) }
func (z *ZapLogger) Fatal(msg string, fields ...zap.Field) { z.l.Fatal(msg, fields...) }
func (z *ZapLogger) Debug(msg string, fields ...zap.Field) { z.l.Debug(msg, fields...) }
func (z *ZapLogger) Warn(msg string, fields ...zap.Field)  { z.l.Warn(msg, fields...) }
