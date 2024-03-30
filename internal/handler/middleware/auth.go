package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	"log"
	"net/http"
)

const (
	sessionTokenKey = "session_token"
	userIdKey       = "user_id"
)

func Auth(l *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {

			if errors.Is(err, http.ErrNoCookie) {
				handler.WriteErrorResponse(w, l, errs.ErrorInfo{
					//GeneralErr: nil,
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
				//GeneralErr: nil,
				LocalErr: errs.ErrUnauthorized,
			})
			return
		}
		if s.(handler.Session).IsExpired() {
			handler.Sessions.Delete(sessionToken)
			handler.WriteErrorResponse(w, l, errs.ErrorInfo{
				//GeneralErr: nil,
				LocalErr: errs.ErrUnauthorized,
			})
			return
		}
		log.Println("OK 1")
		userId := s.(handler.Session).UserId
		log.Println(userId)

		ctx := r.Context()
		log.Println("OK 2")
		//type ctxString string
		//sessionTokenKey := ctxString("session_token")
		//userIdKey := ctxString("user_id")
		log.Println("OK 3")
		ctx = context.WithValue(ctx, sessionTokenKey, sessionToken)
		ctx = context.WithValue(ctx, userIdKey, userId)
		log.Println("OK 4")
		log.Println(ctx.Value("user_id"))

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
			//GeneralErr: nil,
			LocalErr: errs.ErrAlreadyAuthorized,
		})
	}
}
