package service

import (
	"go.uber.org/zap"
	"harmonica/internal/repository"
)

type RepositoryService struct {
	repo *repository.Repository
}

func NewRepositoryService(r *repository.Repository) *RepositoryService {
	return &RepositoryService{repo: r}
}

type Service struct {
	IService
	Logger *zap.Logger
}

func NewService(r *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		IService: NewRepositoryService(r),
		Logger:   logger,
	}
}
