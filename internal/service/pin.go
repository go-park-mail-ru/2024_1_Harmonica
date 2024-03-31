package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, errs.ErrorInfo) {
	pins, err := s.repo.GetFeedPins(ctx, limit, offset)
	if err != nil {
		return entity.FeedPins{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return pins, emptyErrorInfo
}

func (s *RepositoryService) GetUserPins(ctx context.Context, authorId entity.UserID, limit, offset int) (entity.UserPins, errs.ErrorInfo) {
	// TODO check on user exist and throw 404 not 500
	pins, err := s.repo.GetUserPins(ctx, authorId, limit, offset)
	if err != nil {
		return entity.UserPins{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return pins, emptyErrorInfo
}

func (s *RepositoryService) GetPinById(ctx context.Context, pinId entity.PinID) (entity.PinPageResponse, errs.ErrorInfo) {
	pin, err := s.repo.GetPinById(ctx, pinId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return pin, emptyErrorInfo
}

func (s *RepositoryService) CreatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo) {
	pinId, errCreate := s.repo.CreatePin(ctx, pin)
	if errCreate != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: errCreate, // добавила эту ошибку, ранее возвращалось nil
			LocalErr:   errs.ErrDBInternal,
		}
	}
	res, errFind := s.repo.GetPinById(ctx, pinId)
	if errFind != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: errFind,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return res, emptyErrorInfo
}

func (s *RepositoryService) UpdatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo) {
	oldPin, err := s.repo.GetPinById(ctx, pin.PinId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if oldPin.PinAuthor.UserId != pin.AuthorId {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: nil,
			LocalErr:   errs.ErrPermissionDenied,
		}
	}
	err = s.repo.UpdatePin(ctx, pin)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	res, errFind := s.repo.GetPinById(ctx, pin.PinId)
	if errFind != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return res, emptyErrorInfo
}

func (s *RepositoryService) DeletePin(ctx context.Context, pin entity.Pin) errs.ErrorInfo {
	oldPin, err := s.repo.GetPinById(ctx, pin.PinId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if oldPin.PinAuthor.UserId != pin.AuthorId {
		return errs.ErrorInfo{
			GeneralErr: nil,
			LocalErr:   errs.ErrPermissionDenied,
		}
	}
	err = s.repo.DeletePin(ctx, pin.PinId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}
