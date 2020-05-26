package logger

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
