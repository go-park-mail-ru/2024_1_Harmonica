package middleware

import (
	"bufio"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"time"
)

const RequestIdKey = "request_id"

type WrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (wrappedWriter *WrappedResponseWriter) WriteHeader(code int) {
	wrappedWriter.StatusCode = code
	wrappedWriter.ResponseWriter.WriteHeader(code)
}

func (wrappedWriter *WrappedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := wrappedWriter.ResponseWriter.(http.Hijacker)
	if !ok {
		log.Println("hijack not supported")
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
	// это нужно для чата, так как иначе  ошибка "websocket: response does not implement http.Hijacker"
}

func Logging(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestId := uuid.New().String()

		ctx := context.WithValue(r.Context(), RequestIdKey, requestId)
		wrappedWriter := &WrappedResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		l.Info("request handled",
			zap.String("request_id", requestId),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("duration", time.Since(start).String()),
			zap.Int("status_code", wrappedWriter.StatusCode),
		)
	})
}
