package service

import (
	like "harmonica/internal/microservices/like/proto"
	"harmonica/internal/repository"
)

type RepositoryService struct {
	LikeService like.LikeClient
	repo        repository.IRepository
}

func NewRepositoryService(r repository.IRepository, l like.LikeClient) *RepositoryService {
	return &RepositoryService{repo: r, LikeService: l}
}

type Service struct {
	IService
}

func NewService(r repository.IRepository, l like.LikeClient) *Service {
	return &Service{
		IService: NewRepositoryService(r, l),
	}
}

//func NewServiceForTests(s IService) *Service {
//	return &Service{
//		IService: s,
//	}
//}
