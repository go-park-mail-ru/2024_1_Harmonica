package service

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

const ForeignKeyViolationErrCode = pq.ErrorCode("23503")

func (s *RepositoryService) CreateMessage(ctx context.Context, message entity.Message) errs.ErrorInfo {
	message.Sanitize()
	err := s.repo.CreateMessage(ctx, message)
	if err != nil {
		localErr := errs.ErrDBInternal
		var pqErr *pq.Error
		ok := errors.As(err, &pqErr)
		if ok && (pqErr.Code == ForeignKeyViolationErrCode) {
			// юзер, которому пытаемся отправить сообщение, не существует
			localErr = errs.ErrForeignKeyViolation
		}
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   localErr,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) GetMessages(ctx context.Context, firstUserId, secondUserId entity.UserID) (entity.Messages, errs.ErrorInfo) {
	messages, err := s.repo.GetMessages(ctx, firstUserId, secondUserId)
	if err != nil {
		return entity.Messages{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return messages, emptyErrorInfo
}
