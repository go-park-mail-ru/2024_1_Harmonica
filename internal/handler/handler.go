package handler

import (
	//"harmonica/internal/repository"
	"harmonica/internal/service"
)

type APIHandler struct {
	service *service.Service
}

func NewAPIHandler(s *service.Service) *APIHandler {
	return &APIHandler{service: s}
}

//интерфейс между сервисом
