package handler

import (
	"go.uber.org/zap"
	auth "harmonica/internal/microservices/auth/proto"
	"harmonica/internal/service"
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
	service     service.IService
	logger      *zap.Logger
	AuthService auth.AuthorizationClient
}

func NewAPIHandler(s service.IService, l *zap.Logger, a auth.AuthorizationClient) *APIHandler {
	return &APIHandler{
		service:     s,
		logger:      l,
		AuthService: a,
	}
}
