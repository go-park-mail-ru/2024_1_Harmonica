package service

import (
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
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		IService: NewRepositoryService(r),
	}
}
