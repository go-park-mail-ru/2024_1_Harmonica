package service

import (
	"harmonica/internal/entity"
	"harmonica/internal/repository"
	"log"
)

type PinService struct {
	repo *repository.Repository
}

func NewPinService(r *repository.Repository) *PinService {
	return &PinService{repo: r}
}

func (p PinService) GetPins(limit, offset int) (entity.Pins, error) {
	user, _ := p.repo.GetUserById(18)
	log.Println(user.Email)
	pins, err := p.repo.GetPins(limit, offset)
	if err != nil {
		return pins, err
	}
	return pins, nil
}
