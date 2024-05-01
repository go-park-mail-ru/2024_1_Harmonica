package handler

import (
	"go.uber.org/zap"
	auth "harmonica/internal/microservices/auth/proto"
	image "harmonica/internal/microservices/image/proto"
	"harmonica/internal/service"
)

type APIHandler struct {
	service      service.IService
	logger       *zap.Logger
	hub          *Hub
	AuthService  auth.AuthorizationClient
	ImageService image.ImageClient
}

func NewAPIHandler(s service.IService, l *zap.Logger, hub *Hub, a auth.AuthorizationClient, i image.ImageClient) *APIHandler {
	return &APIHandler{
		service:      s,
		logger:       l,
		hub:          hub,
		AuthService:  a,
		ImageService: i,
	}
}
