package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

type IService interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, errs.ErrorInfo)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, errs.ErrorInfo)
	GetUserProfileByNickname(ctx context.Context, nickname string, userId entity.UserID) (entity.UserProfileResponse, errs.ErrorInfo)
	GetUserById(ctx context.Context, id entity.UserID) (entity.User, errs.ErrorInfo)
	RegisterUser(ctx context.Context, user entity.User) []errs.ErrorInfo
	UpdateUser(ctx context.Context, user entity.User) (entity.User, errs.ErrorInfo)

	GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, errs.ErrorInfo)
	GetSubscriptionsFeedPins(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, errs.ErrorInfo)
	GetUserPins(ctx context.Context, authorNickname string, limit, offset int) (entity.UserPins, errs.ErrorInfo)
	GetPinById(ctx context.Context, PinId entity.PinID, UserId entity.UserID) (entity.PinPageResponse, errs.ErrorInfo)
	CreatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo)
	UpdatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo)
	DeletePin(ctx context.Context, pin entity.Pin) errs.ErrorInfo

	CreateBoard(ctx context.Context, board entity.Board, userId entity.UserID) (entity.FullBoard, errs.ErrorInfo)
	GetBoardById(ctx context.Context, boardId entity.BoardID, userId entity.UserID, limit, offset int) (entity.FullBoard, errs.ErrorInfo)
	UpdateBoard(ctx context.Context, board entity.Board, userId entity.UserID) (entity.FullBoard, errs.ErrorInfo)
	AddPinToBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo
	DeletePinFromBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo
	DeleteBoard(ctx context.Context, boardId entity.BoardID, userId entity.UserID) errs.ErrorInfo
	GetUserBoards(ctx context.Context, authorNickname string, userId entity.UserID, limit, offset int) (entity.UserBoards, errs.ErrorInfo)

	CreateMessage(ctx context.Context, message entity.Message) errs.ErrorInfo
	GetMessages(ctx context.Context, dialogUserId, authUserId entity.UserID) (entity.Messages, errs.ErrorInfo)
	GetUserChats(ctx context.Context, userId entity.UserID) (entity.UserChats, errs.ErrorInfo)

	UpdateDraft(ctx context.Context, draft entity.Draft) errs.ErrorInfo

	AddSubscriptionToUser(ctx context.Context, userId, subscribeUserId entity.UserID) errs.ErrorInfo
	DeleteSubscriptionToUser(ctx context.Context, userId, unsubscribeUserId entity.UserID) errs.ErrorInfo
	GetUserSubscribers(ctx context.Context, userId entity.UserID) (entity.UserSubscribers, errs.ErrorInfo)
	GetUserSubscriptions(ctx context.Context, userId entity.UserID) (entity.UserSubscriptions, errs.ErrorInfo)

	Search(ctx context.Context, query string) (entity.SearchResult, errs.ErrorInfo)
}
