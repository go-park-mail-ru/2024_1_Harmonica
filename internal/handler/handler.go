package handler

import (
	auth "harmonica/internal/microservices/auth/proto"
	image "harmonica/internal/microservices/image/proto"
	like "harmonica/internal/microservices/like/proto"
	"harmonica/internal/service"

	"go.uber.org/zap"
)

type APIHandler struct {
	service      service.IService
	logger       *zap.Logger
	hub          *Hub
	AuthService  auth.AuthorizationClient
	ImageService image.ImageClient
	LikeService  like.LikeClient
}

func NewAPIHandler(s service.IService, l *zap.Logger, hub *Hub, a auth.AuthorizationClient,
	i image.ImageClient, like like.LikeClient) *APIHandler {
	return &APIHandler{
		service:      s,
		logger:       l,
		hub:          hub,
		AuthService:  a,
		ImageService: i,
		LikeService:  like,
	}
}
