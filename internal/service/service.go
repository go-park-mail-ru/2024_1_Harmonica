package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/repository"
)

type IService interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, error)
	GetUserById(ctx context.Context, id int64) (entity.User, error)
	RegisterUser(ctx context.Context, user entity.User) []error
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	GetPins(ctx context.Context, limit, offset int) (entity.Pins, error)
}

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

//type User interface {
//	GetUserByEmail(email string) (entity.User, error)
//	GetUserByNickname(nickname string) (entity.User, error)
//	GetUserById(id int64) (entity.User, error)
//	RegisterUser(user entity.User) []error
//}
//
//type Pin interface {
//	GetPins(limit, offset int) (entity.Pins, error)
//}

//func NewService(r *repository.Repository) *Service {
//	return &Service{
//		User: NewUserService(r),
//		Pin:  NewPinService(r),
//	}
//}
