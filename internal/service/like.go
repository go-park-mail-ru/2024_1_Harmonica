package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (r *RepositoryService) SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	err := r.repo.SetLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrDBInternal
	}
	return nil
}

func (r *RepositoryService) ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	err := r.repo.ClearLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrDBInternal
	}
	return nil
}

func (r *RepositoryService) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error) {
	res, err := r.repo.GetUsersLiked(ctx, pinId, limit)
	if err != nil {
		return entity.UserList{}, errs.ErrDBInternal
	}
	return res, nil
}
