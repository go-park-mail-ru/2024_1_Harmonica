package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cors.New(cors.Options{
			AllowedOrigins:     []string{"http://localhost:8000", "http://85.192.35.36:8000"},
			AllowCredentials:   true,
			AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:     []string{"*"},
			OptionsPassthrough: false,
		})
		c.Handler(next)
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
