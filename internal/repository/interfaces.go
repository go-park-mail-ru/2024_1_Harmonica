package repository

import (
	"context"
	"harmonica/internal/entity"
)

type IRepository interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, error)
	GetUserById(ctx context.Context, id entity.UserID) (entity.User, error)
	RegisterUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, user entity.User) error

	GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, error)
	GetSubscriptionsFeedPins(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, error)
	GetUserPins(ctx context.Context, authorId entity.UserID, limit, offset int) (entity.UserPins, error)
	GetPinById(ctx context.Context, PinId entity.PinID) (entity.PinPageResponse, error)
	CreatePin(ctx context.Context, pin entity.Pin) (entity.PinID, error)
	UpdatePin(ctx context.Context, pin entity.Pin) error
	DeletePin(ctx context.Context, id entity.PinID) error
	CheckPinExistence(ctx context.Context, id entity.PinID) (bool, error)

	CreateBoard(ctx context.Context, board entity.Board, userId entity.UserID) (entity.Board, error)
	GetBoardById(ctx context.Context, boardId entity.BoardID) (entity.Board, error)
	GetBoardAuthors(ctx context.Context, boardId entity.BoardID) ([]entity.BoardAuthor, error)
	GetBoardPins(ctx context.Context, boardId entity.BoardID, limit, offset int) ([]entity.BoardPinResponse, error)
	UpdateBoard(ctx context.Context, board entity.Board) (entity.Board, error)
	AddPinToBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID) error
	DeletePinFromBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID) error
	DeleteBoard(ctx context.Context, boardId entity.BoardID) error
	GetUserBoards(ctx context.Context, authorId, userId entity.UserID, limit, offset int) (entity.UserBoards, error)
	CheckBoardAuthorExistence(ctx context.Context, userId entity.UserID, boardId entity.BoardID) (bool, error)

	CreateMessage(ctx context.Context, message entity.Message) error
	GetMessages(ctx context.Context, firstUserId, secondUserId entity.UserID) (entity.Messages, error)
	GetUserChats(ctx context.Context, userId entity.UserID) (entity.UserChats, error)

	AddSubscriptionToUser(ctx context.Context, userId, subscribeUserId entity.UserID) error
	DeleteSubscriptionToUser(ctx context.Context, userId, unsubscribeUserId entity.UserID) error
	GetSubscriptionsInfo(ctx context.Context, userToGetInfoId, userId entity.UserID) (entity.UserProfileResponse, error)
	GetUserSubscribers(ctx context.Context, userId entity.UserID) (entity.UserSubscribers, error)
	GetUserSubscriptions(ctx context.Context, userId entity.UserID) (entity.UserSubscriptions, error)

	SearchForUsers(ctx context.Context, query string) ([]entity.SearchUser, error)
	SearchForPins(ctx context.Context, query string) ([]entity.SearchPin, error)
	SearchForBoards(ctx context.Context, query string) ([]entity.SearchBoard, error)
}
