package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) AddComment(ctx context.Context, comment string, pinId entity.PinID, userId entity.UserID) (entity.PinPageResponse, errs.ErrorInfo) {
	exists, err := s.repo.CheckPinExistence(ctx, pinId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !exists {
		return entity.PinPageResponse{}, errs.ErrorInfo{LocalErr: errs.ErrNotFound}
	}
	if len(comment) <= 0 {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			LocalErr: errs.ErrEmptyComment,
		}
	}
	commentId, err := s.repo.AddComment(ctx, comment, pinId, userId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	pin, err := s.repo.GetPinById(ctx, pinId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	n := entity.Notification{
		Type:              entity.NotificationTypeComment,
		UserId:            pin.PinAuthor.UserId,
		TriggeredByUserId: userId,
		CommentId:         commentId,
		PinId:             pinId,
	}
	err = s.repo.CreateNotification(ctx, n)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return pin, errs.ErrorInfo{}
}

func (s *RepositoryService) GetComments(ctx context.Context, pinId entity.PinID) (entity.GetCommentsResponse, errs.ErrorInfo) {
	exists, err := s.repo.CheckPinExistence(ctx, pinId)
	if err != nil {
		return entity.GetCommentsResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !exists {
		return entity.GetCommentsResponse{}, errs.ErrorInfo{LocalErr: errs.ErrNotFound}
	}
	res, err := s.repo.GetComments(ctx, pinId)
	if err != nil {
		return entity.GetCommentsResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return res, errs.ErrorInfo{}
}
