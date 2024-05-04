package service

import "harmonica/internal/microservices/image/server/repository"

type RepositoryService struct {
	repo repository.IRepository
}

func NewRepositoryService(r repository.IRepository) *RepositoryService {
	return &RepositoryService{repo: r}
}

type Service struct {
	IService
}

func NewService(r repository.IRepository) *Service {
	return &Service{
		IService: NewRepositoryService(r),
	}
}
