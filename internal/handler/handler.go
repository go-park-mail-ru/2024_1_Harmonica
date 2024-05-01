package handler

import (
	auth "harmonica/internal/microservices/auth/proto"
	image "harmonica/internal/microservices/image/proto"
	"harmonica/internal/service"

	"go.uber.org/zap"
)

//type APIHandler struct {
//	service *service.Service
//	logger  *zap.Logger
//}
//
//func NewAPIHandler(s *service.Service, l *zap.Logger) *APIHandler {
//	return &APIHandler{
//		service: s,
//		logger:  l,
//	}
//}

// тут ощутимые изменения. не накосячила?

type APIHandler struct {
	service      service.IService
	logger       *zap.Logger
	AuthService  auth.AuthorizationClient
	ImageService image.ImageClient
}

func NewAPIHandler(s service.IService, l *zap.Logger, a auth.AuthorizationClient, i image.ImageClient) *APIHandler {
	return &APIHandler{
		service:      s,
		logger:       l,
		AuthService:  a,
		ImageService: i,
	}
}
