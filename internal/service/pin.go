package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	like "harmonica/internal/microservices/like/proto"
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

func (s *RepositoryService) GetSubscriptionsFeedPins(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, errs.ErrorInfo) {
	pins, err := s.repo.GetSubscriptionsFeedPins(ctx, userId, limit, offset)
	if err != nil {
		return entity.FeedPins{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return pins, emptyErrorInfo
}

func (s *RepositoryService) GetUserPins(ctx context.Context, authorNickname string, limit, offset int) (entity.UserPins, errs.ErrorInfo) {
	user, err := s.repo.GetUserByNickname(ctx, authorNickname)
	if err != nil {
		return entity.UserPins{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	pins, errPin := s.repo.GetUserPins(ctx, user.UserID, limit, offset)
	if errPin != nil {
		return entity.UserPins{}, errs.ErrorInfo{
			GeneralErr: errPin,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return pins, emptyErrorInfo
}

func (s *RepositoryService) GetPinById(ctx context.Context, pinId entity.PinID, userId entity.UserID) (entity.PinPageResponse, errs.ErrorInfo) {
	pin, err := s.repo.GetPinById(ctx, pinId)
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrElementNotExist,
		}
	}
	res, err := s.LikeService.CheckIsLiked(ctx, &like.CheckIsLikedRequest{PinId: int64(pinId), UserId: int64(userId)})
	if err != nil {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			LocalErr: errs.ErrGRPCWentWrong,
		}
	}
	if !res.Valid {
		return entity.PinPageResponse{}, errs.ErrorInfo{
			LocalErr: errs.GetLocalErrorByCode[res.LocalError],
		}
	}
	pin.IsLiked = res.Liked
	return pin, emptyErrorInfo
}

func (s *RepositoryService) CreatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo) {
	pin.Sanitize()
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
			LocalErr:   errs.ErrElementNotExist,
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
			GeneralErr: errFind,
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
			LocalErr:   errs.ErrElementNotExist,
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
