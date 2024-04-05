package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	"net/http"
)

const (
	SessionTokenKey = "session_token"
	UserIdKey       = "user_id"
	IsAuthKey       = "is_auth"
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
		ctx = context.WithValue(ctx, SessionTokenKey, sessionToken)
		ctx = context.WithValue(ctx, UserIdKey, userId)
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
		if exists && !s.(handler.Session).IsExpired() {
			handler.WriteErrorResponse(w, l, errs.ErrorInfo{
				LocalErr: errs.ErrAlreadyAuthorized,
			})
			return
		}

		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		if s.(handler.Session).IsExpired() {
			handler.Sessions.Delete(sessionToken)
			next.ServeHTTP(w, r)
			return
		}
	}
}

func CheckAuth(l *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		ctx := r.Context()
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				ctx = context.WithValue(ctx, IsAuthKey, false)
				next.ServeHTTP(w, r.WithContext(ctx))
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
			ctx = context.WithValue(ctx, IsAuthKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if s.(handler.Session).IsExpired() {
			handler.Sessions.Delete(sessionToken)
			ctx = context.WithValue(ctx, IsAuthKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userId := s.(handler.Session).UserId
		ctx = context.WithValue(ctx, IsAuthKey, true)
		ctx = context.WithValue(ctx, SessionTokenKey, sessionToken)
		ctx = context.WithValue(ctx, UserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
