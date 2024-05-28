package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) AddComment(ctx context.Context, comment string, pinId entity.PinID, userId entity.UserID) (entity.PinPageResponse, entity.CommentID, errs.ErrorInfo) {
	exists, err := s.repo.CheckPinExistence(ctx, pinId)
	if err != nil {
		return entity.PinPageResponse{}, entity.CommentID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if !exists {
		return entity.PinPageResponse{}, entity.CommentID(0), errs.ErrorInfo{LocalErr: errs.ErrNotFound}
	}
	if len(comment) <= 0 {
		return entity.PinPageResponse{}, entity.CommentID(0), errs.ErrorInfo{LocalErr: errs.ErrEmptyComment}
	}
	commentId, err := s.repo.AddComment(ctx, comment, pinId, userId)
	if err != nil {
		return entity.PinPageResponse{}, entity.CommentID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	// пин, к которому комментарий, чтобы позже передать инфу о нем в уведомление
	pin, err := s.repo.GetPinById(ctx, pinId)
	if err != nil {
		return entity.PinPageResponse{}, entity.CommentID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	//n := entity.Notification{
	//	Type:              entity.NotificationTypeComment,
	//	UserId:            pin.PinAuthor.UserId,
	//	TriggeredByUserId: userId,
	//	CommentId:         commentId,
	//	PinId:             pinId,
	//}
	//_, err = s.repo.CreateNotification(ctx, n)
	//if err != nil {
	//	return entity.PinPageResponse{}, entity.CommentID(0), errs.ErrorInfo{
	//		GeneralErr: err,
	//		LocalErr:   errs.ErrDBInternal,
	//	}
	//}
	return pin, commentId, errs.ErrorInfo{}
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
