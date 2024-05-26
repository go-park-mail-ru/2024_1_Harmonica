package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) GetUnreadNotifications(ctx context.Context, userId entity.UserID) (entity.Notifications, errs.ErrorInfo) {
	notifications, err := s.repo.GetUnreadNotifications(ctx, userId)
	if err != nil {
		return entity.Notifications{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return notifications, emptyErrorInfo
}

func (s *RepositoryService) CreateNotification(ctx context.Context, n entity.Notification) errs.ErrorInfo {
	err := s.repo.CreateNotification(ctx, n)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}
