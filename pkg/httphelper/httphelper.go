package httphelper

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	WriteJSON(w, status, Success(data))
}

func WriteSuccessWithMeta(w http.ResponseWriter, status int, data, meta interface{}) {
	WriteJSON(w, status, SuccessWithMeta(data, meta))
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, Error(msg))
}

// ParseInt — безопасное преобразование строки в число
func ParseInt(s string, def int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
