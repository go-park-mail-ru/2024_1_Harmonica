package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (r *RepositoryService) GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, error) {
	pins, err := r.repo.GetFeedPins(ctx, limit, offset)
	if err != nil {
		return entity.FeedPins{}, errs.ErrDBInternal
	}
	return pins, nil
}

func (r *RepositoryService) GetUserPins(ctx context.Context, authorId entity.UserID, limit, offset int) (entity.UserPins, error) {
	// TODO check on user exist and throw 404 not 500
	pins, err := r.repo.GetUserPins(ctx, authorId, limit, offset)
	if err != nil {
		return entity.UserPins{}, errs.ErrDBInternal
	}
	return pins, nil
}

func (r *RepositoryService) GetPinById(ctx context.Context, pinId entity.PinID) (entity.PinPageResponse, error) {
	pin, err := r.repo.GetPinById(ctx, pinId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrDBInternal
	}
	return pin, nil
}

func (r *RepositoryService) CreatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, error) {
	pinId, errCreate := r.repo.CreatePin(ctx, pin)
	if errCreate != nil {
		return entity.PinPageResponse{}, nil
	}
	res, errFind := r.repo.GetPinById(ctx, pinId)
	if errFind != nil {
		return entity.PinPageResponse{}, errs.ErrDBInternal
	}
	return res, nil
}

func (r *RepositoryService) UpdatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, error) {
	oldPin, err := r.repo.GetPinById(ctx, pin.PinId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrDBInternal
	}
	if oldPin.PinAuthor.UserId != pin.AuthorId {
		return entity.PinPageResponse{}, errs.ErrPermissionDenied
	}
	err = r.repo.UpdatePin(ctx, pin)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrDBInternal
	}
	res, errFind := r.repo.GetPinById(ctx, pin.PinId)
	if errFind != nil {
		return entity.PinPageResponse{}, errs.ErrDBInternal
	}
	return res, nil
}

func (r *RepositoryService) DeletePin(ctx context.Context, pin entity.Pin) error {
	oldPin, err := r.repo.GetPinById(ctx, pin.PinId)
	if err != nil {
		return errs.ErrDBInternal
	}
	if oldPin.PinAuthor.UserId != pin.AuthorId {
		return errs.ErrPermissionDenied
	}
	err = r.repo.DeletePin(ctx, pin.PinId)
	if err != nil {
		return errs.ErrDBInternal
	}
	return nil
}
