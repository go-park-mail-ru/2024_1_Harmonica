package repository

import (
	"context"
	"harmonica/internal/entity"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type IRepository interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (entity.User, error)
	GetUserById(ctx context.Context, id entity.UserID) (entity.User, error)
	RegisterUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, user entity.User) error

	GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, error)
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
	GetUserBoards(ctx context.Context, authorId entity.UserID, limit, offset int) (entity.UserBoards, error)
	CheckBoardAuthorExistence(ctx context.Context, userId entity.UserID, boardId entity.BoardID) (bool, error)

	SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error
	ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error
	GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error)
	CheckIsLiked(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error)

	UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	GetImage(ctx context.Context, name string) (*minio.Object, error)

	CreateMessage(ctx context.Context, message entity.Message) error
	GetMessages(ctx context.Context, firstUserId, secondUserId entity.UserID) (entity.Messages, error)
	GetUserChats(ctx context.Context, userId entity.UserID) (entity.UserChats, error)
}
