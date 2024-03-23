package handler

import (
	"context"
	"fmt"
	"github.com/rs/cors"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		sessionToken := c.Value
		s, exists := sessions.Load(sessionToken)
		if !exists {
			fmt.Println("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		if s.(Session).IsExpired() {
			sessions.Delete(sessionToken)
			fmt.Println("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		//return sessionToken, s.(Session).UserId, nil

		ctx := r.Context()
		ctx = context.WithValue(ctx, "session_token", sessionToken)
		ctx = context.WithValue(ctx, "user_id", s.(Session).UserId)

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cors.New(cors.Options{
			AllowedOrigins:     []string{"http://localhost:8000", "http://85.192.35.36:8000"},
			AllowCredentials:   true,
			AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:     []string{"*"},
			OptionsPassthrough: false,
		})
		c.Handler(next)
		next.ServeHTTP(w, r)
	})
}
