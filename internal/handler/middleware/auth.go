package middleware

import (
	"context"
	"errors"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	"net/http"

	"go.uber.org/zap"
)

const (
	SessionTokenKey = "session_token"
	UserIdKey       = "user_id"
	IsAuthKey       = "is_auth"
)

func CheckSession(r *http.Request) (*http.Request, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	sessionToken := c.Value
	s, exists := handler.Sessions.Load(sessionToken)
	if !exists || s.(handler.Session).IsExpired() {
		if exists {
			handler.Sessions.Delete(sessionToken)
		}
		return nil, errs.ErrUnauthorized
	}
	userId := s.(handler.Session).UserId
	ctx := r.Context()
	ctx = context.WithValue(ctx, UserIdKey, userId)
	return r.WithContext(ctx), nil
}

func AuthRequired(l *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := CheckSession(r)
		if err != nil {
			if errs.ErrorCodes[err].HttpCode != 0 {
				handler.WriteErrorResponse(w, l, errs.ErrorInfo{LocalErr: err})
				return
			}
			if errors.Is(err, http.ErrNoCookie) {
				handler.WriteErrorResponse(w, l, errs.ErrorInfo{LocalErr: errs.ErrUnauthorized})
			}
			handler.WriteErrorResponse(w, l, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrReadCookie})
			return
		}
		next.ServeHTTP(w, request)
	}
}

func NoAuthRequired(l *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckSession(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		handler.WriteErrorResponse(w, l, errs.ErrorInfo{
			LocalErr: errs.ErrAlreadyAuthorized,
		})
	}
}

func CheckAuth(l *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := CheckSession(r)
		if err != nil {
			request = r
		}
		next.ServeHTTP(w, request)
	}
}
