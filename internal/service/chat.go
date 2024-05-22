package service

import (
	"context"
	"database/sql"
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
	err = s.repo.UpdateDraft(ctx, entity.Draft{SenderId: message.SenderId, ReceiverId: message.ReceiverId, Text: ""})
	if err != nil {
		return errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return emptyErrorInfo
}

func (s *RepositoryService) GetMessages(ctx context.Context, dialogUserId, authUserId entity.UserID) (entity.Messages, errs.ErrorInfo) {
	messages, err := s.repo.GetMessages(ctx, dialogUserId, authUserId)
	if err != nil {
		return entity.Messages{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	messages.Draft, err = s.repo.GetDraft(ctx, dialogUserId, authUserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entity.Messages{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return messages, emptyErrorInfo
}

func (s *RepositoryService) GetUserChats(ctx context.Context, userId entity.UserID) (entity.UserChats, errs.ErrorInfo) {
	chats, err := s.repo.GetUserChats(ctx, userId)
	if err != nil {
		return entity.UserChats{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return chats, emptyErrorInfo
}
