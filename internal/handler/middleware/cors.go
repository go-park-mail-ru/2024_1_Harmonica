package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func CORS(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"https://harmoniums.ru", "http://localhost", "https://localhost"},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})
	return c.Handler(next)
}
