package middlewares

import (
	"dozenChairs/pkg/httphelper"
	"dozenChairs/pkg/logger"
	"net/http"
	"runtime/debug"
)

func Recover(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						logger.ZapAny("recover", rec),
						logger.ZapAny("stack", string(debug.Stack())),
					)
					httphelper.WriteError(w, http.StatusInternalServerError, "Internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
