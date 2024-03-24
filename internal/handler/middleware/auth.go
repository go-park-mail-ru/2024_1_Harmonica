package middleware

import (
	"context"
	"errors"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	"net/http"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")
		if err != nil {

			if errors.Is(err, http.ErrNoCookie) {
				handler.WriteErrorResponse(w, errs.ErrUnauthorized)
				return
			}

			handler.WriteErrorResponse(w, errs.ErrReadCookie)
			return
		}

		sessionToken := c.Value
		s, exists := handler.Sessions.Load(sessionToken)
		if !exists {
			handler.WriteErrorResponse(w, errs.ErrUnauthorized)
			return
		}
		if s.(handler.Session).IsExpired() {
			handler.Sessions.Delete(sessionToken)
			handler.WriteErrorResponse(w, errs.ErrUnauthorized)
			return
		}

		ctx := r.Context()
		type ctxString string
		sessionTokenKey := ctxString("session_token")
		userIdKey := ctxString("user_id")
		ctx = context.WithValue(ctx, sessionTokenKey, sessionToken)
		ctx = context.WithValue(ctx, userIdKey, s.(handler.Session).UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func NotAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session_token")
		if err != nil {

			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
				return
			}

			handler.WriteErrorResponse(w, errs.ErrReadCookie)
			return
		}

		sessionToken := c.Value
		s, exists := handler.Sessions.Load(sessionToken)
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		if s.(handler.Session).IsExpired() {
			handler.Sessions.Delete(sessionToken)
			next.ServeHTTP(w, r)
			return
		}

		handler.WriteErrorResponse(w, errs.ErrAlreadyAuthorized)
	}
}
