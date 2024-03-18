package service

import (
	"harmonica/internal/entity"
	"harmonica/internal/repository"
)

type User interface {
	GetUserByEmail(email string) (entity.User, error)
	GetUserByNickname(nickname string) (entity.User, error)
	GetUserById(id int64) (entity.User, error)
	RegisterUser(user entity.User) []error
}

type Pin interface {
	GetPins(limit, offset int) (entity.Pins, error)
}

type Service struct {
	User
	Pin
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		User: NewUserService(r),
		Pin:  NewPinService(r),
	}
}
