package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func CORS(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"http://localhost:8000", "http://85.192.35.36:8000"},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})
	return c.Handler(next)
}
