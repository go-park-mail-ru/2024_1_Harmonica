package service

import (
	"context"
	"database/sql"
	"errors"
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

func (s *RepositoryService) CreateNotification(ctx context.Context, n entity.Notification) (entity.NotificationID, errs.ErrorInfo) {
	notificationId, err := s.repo.CreateNotification(ctx, n)
	if err != nil {
		return entity.NotificationID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return notificationId, emptyErrorInfo
}

func (s *RepositoryService) ReadNotification(ctx context.Context, notificationId entity.NotificationID, userId entity.UserID) errs.ErrorInfo {
	notification, err := s.repo.GetNotificationById(ctx, notificationId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	if notification.UserId != userId || errors.Is(err, sql.ErrNoRows) {
		// если уведомление принадлежит не тому юзеру, от которого пришел запрос на прочтение, ничего не делаем
		return emptyErrorInfo
	}
	err = s.repo.ReadNotification(ctx, notificationId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) ReadAllNotifications(ctx context.Context, userId entity.UserID) errs.ErrorInfo {
	err := s.repo.ReadAllNotifications(ctx, userId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}
