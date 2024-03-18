package service

import (
	"harmonica/internal/entity"
)

//type PinService struct {
//	repo *repository.Repository
//}
//
//func NewPinService(r *repository.Repository) *PinService {
//	return &PinService{repo: r}
//}

func (r RepositoryService) GetPins(limit, offset int) (entity.Pins, error) {
	pins, err := r.repo.GetPins(limit, offset)
	if err != nil {
		return pins, err
	}
	return pins, nil
}
