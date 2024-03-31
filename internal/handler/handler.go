package handler

import (
	//"harmonica/internal/repository"
	"go.uber.org/zap"
	"harmonica/internal/service"
)

type APIHandler struct {
	service *service.Service
	logger  *zap.Logger
}

func NewAPIHandler(s *service.Service, l *zap.Logger) *APIHandler {
	return &APIHandler{
		service: s,
		logger:  l,
	}
}
