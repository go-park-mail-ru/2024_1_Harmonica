package service

import (
	"harmonica/internal/entity"
	"harmonica/internal/repository"
)

type IService interface {
	GetUserByEmail(email string) (entity.User, error)
	GetUserByNickname(nickname string) (entity.User, error)
	GetUserById(id int64) (entity.User, error)
	RegisterUser(user entity.User) []error
	GetPins(limit, offset int) (entity.Pins, error)
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
