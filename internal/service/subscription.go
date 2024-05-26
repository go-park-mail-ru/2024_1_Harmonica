package service

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

func (s *RepositoryService) AddSubscriptionToUser(ctx context.Context,
	userId, subscribeUserId entity.UserID) errs.ErrorInfo {
	err := s.repo.AddSubscriptionToUser(ctx, userId, subscribeUserId)
	if err != nil {
		if err != nil {
			localErr := errs.ErrDBInternal
			var pqErr *pq.Error
			ok := errors.As(err, &pqErr)
			if ok && (pqErr.Code == UniqueViolationErrCode) {
				localErr = errs.ErrDBUniqueViolation
			}
			if ok && (pqErr.Code == ForeignKeyViolationErrCode) {
				localErr = errs.ErrForeignKeyViolation
			}
			return errs.ErrorInfo{
				GeneralErr: err,
				LocalErr:   localErr,
			}
		}
	}
	n := entity.Notification{
		Type:              entity.NotificationTypeSubscription,
		UserId:            subscribeUserId,
		TriggeredByUserId: userId,
	}
	err = s.repo.CreateNotification(ctx, n)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) DeleteSubscriptionToUser(ctx context.Context,
	userId, unsubscribeUserId entity.UserID) errs.ErrorInfo {
	err := s.repo.DeleteSubscriptionToUser(ctx, userId, unsubscribeUserId)
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) GetUserSubscribers(ctx context.Context, userId entity.UserID) (entity.UserSubscribers, errs.ErrorInfo) {
	subscribers, err := s.repo.GetUserSubscribers(ctx, userId)
	if err != nil {
		return entity.UserSubscribers{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return subscribers, emptyErrorInfo
}

func (s *RepositoryService) GetUserSubscriptions(ctx context.Context, userId entity.UserID) (entity.UserSubscriptions, errs.ErrorInfo) {
	subscriptions, err := s.repo.GetUserSubscriptions(ctx, userId)
	if err != nil {
		return entity.UserSubscriptions{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return subscriptions, emptyErrorInfo
}
