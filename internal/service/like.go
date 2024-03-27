package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (r *RepositoryService) CreateLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	exist, err := r.repo.FindLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrDBInternal
	}
	if exist {
		return errs.ErrLikeAlreadyCreated
	}
	err = r.repo.CreateLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrDBInternal
	}
	return nil
}

func (r *RepositoryService) DeleteLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	exist, err := r.repo.FindLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrDBInternal
	}
	if !exist {
		return errs.ErrLikeAlreadyDeleted
	}
	err = r.repo.DeleteLike(ctx, pinId, userId)
	if err != nil {
		return errs.ErrDBInternal
	}
	return nil
}

func (r *RepositoryService) GetLikedPins(ctx context.Context, userId entity.UserID, limit int) (entity.FeedPins, error) {
	res, err := r.repo.GetLikedPins(ctx, userId, limit)
	if err != nil {
		return entity.FeedPins{}, errs.ErrDBInternal
	}
	return res, nil
}

func (r *RepositoryService) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error) {
	res, err := r.repo.GetUsersLiked(ctx, pinId, limit)
	if err != nil {
		return entity.UserList{}, errs.ErrDBInternal
	}
	return res, nil
}

func (r *RepositoryService) GetLikesCount(ctx context.Context, pinId entity.PinID) (uint64, error) {
	res, err := r.repo.GetLikesCount(ctx, pinId)
	if err != nil {
		return 0, errs.ErrDBInternal
	}
	return res, nil
}
