package service

import (
	"context"
	"harmonica/internal/entity"
)

func (r *RepositoryService) GetPins(ctx context.Context, limit, offset int) (entity.Pins, error) {
	pins, err := r.repo.GetPins(ctx, limit, offset)
	if err != nil {
		return pins, err
	}
	return pins, nil
}
