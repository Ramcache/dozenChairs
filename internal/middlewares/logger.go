package middlewares

import (
	"dozenChairs/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: 200}

			next.ServeHTTP(rw, r)

			log.Info("http request",
				logger.ZapStr("method", r.Method),
				logger.ZapStr("path", r.URL.Path),
				logger.ZapInt("status", rw.status),
				logger.ZapStr("remote", r.RemoteAddr),
				logger.ZapStr("request_id", GetRequestID(r.Context())),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
