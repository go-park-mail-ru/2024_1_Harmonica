package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	"net/http"
)

type contextKey string

const (
	sessionTokenKey contextKey = "session_token"
	userIdKey       contextKey = "user_id"
)

func Auth(l *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {

			if errors.Is(err, http.ErrNoCookie) {
				handler.WriteErrorResponse(w, l, errs.ErrorInfo{
					LocalErr: errs.ErrUnauthorized,
				})
				return
			}

			handler.WriteErrorResponse(w, l, errs.ErrorInfo{
				GeneralErr: err,
				LocalErr:   errs.ErrReadCookie,
			})
			return
		}

		sessionToken := c.Value
		s, exists := handler.Sessions.Load(sessionToken)
		if !exists {
			handler.WriteErrorResponse(w, l, errs.ErrorInfo{
				LocalErr: errs.ErrUnauthorized,
			})
			return
		}
		if s.(handler.Session).IsExpired() {
			handler.Sessions.Delete(sessionToken)
			handler.WriteErrorResponse(w, l, errs.ErrorInfo{
				LocalErr: errs.ErrUnauthorized,
			})
			return
		}

		ctx := r.Context()
		userId := s.(handler.Session).UserId
		ctx = context.WithValue(ctx, sessionTokenKey, sessionToken)
		ctx = context.WithValue(ctx, userIdKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func NotAuth(l *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {

			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
				return
			}

			handler.WriteErrorResponse(w, l, errs.ErrorInfo{
				GeneralErr: err,
				LocalErr:   errs.ErrReadCookie,
			})
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

		handler.WriteErrorResponse(w, l, errs.ErrorInfo{
			LocalErr: errs.ErrAlreadyAuthorized,
		})
	}
}
