package handlers

import (
	"log"
	"net/http"
)

// LoggingMiddleware registra el método, la URL y la IP de cada petición
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s (desde %s)", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
