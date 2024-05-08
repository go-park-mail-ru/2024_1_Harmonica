package server

import (
	"context"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	like "harmonica/internal/microservices/like/proto"
	"harmonica/internal/microservices/like/server/service"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var (
	emptyErrorInfo      = errs.ErrorInfo{}
	SetLikeMethod       = `SetLike(%d, %d)`
	ClearLikeMethod     = `ClearLike(%d, %d)`
	GetUsersLikedMethod = `GetUsersLiked(%d, %d)`
	CheckIsLikedMethod  = `CheckIsLiked(%d, %d)`
	GetFavoritesMethod  = `GetFavorites(%d, %d, %d)`
)

type LikeServer struct {
	service service.IService
	logger  *zap.Logger
	like.LikeServer
}

func NewLikeServer(s *service.Service, l *zap.Logger) LikeServer {
	return LikeServer{service: s, logger: l}
}

func GetRequestId(ctx context.Context) string {
	if len(metadata.ValueFromIncomingContext(ctx, "request_id")) > 0 {
		return metadata.ValueFromIncomingContext(ctx, "request_id")[0]
	}
	return ""
}

func (s LikeServer) logErr(method string, err error, requestId string) {
	s.logger.Error(fmt.Sprintf("Error occurred in '%s'", method),
		zap.String("request_id", requestId),
		zap.Error(err))
}

func (s LikeServer) logHandle(request string, requestId string) {
	s.logger.Info(fmt.Sprintf("Request handled: '%s'", request),
		zap.String("request_id", requestId),
	)
}

func MakeRequestString(method string, args ...any) string {
	return fmt.Sprintf(method, args...)
}

func (s LikeServer) SetLike(ctx context.Context, req *like.MakeLikeRequest) (*like.MakeLikeResponse, error) {
	requestId := GetRequestId(ctx)
	requestString := MakeRequestString(SetLikeMethod, req.PinId, req.UserId)
	s.logHandle(requestString, requestId)

	err := s.service.SetLike(ctx, entity.PinID(req.PinId), entity.UserID(req.UserId))
	if err != emptyErrorInfo {
		s.logErr(requestString, err.GeneralErr, requestId)
		return &like.MakeLikeResponse{Valid: false, LocalError: int64(errs.ErrorCodes[err.LocalErr].LocalCode)}, nil
	}
	return &like.MakeLikeResponse{Valid: true}, nil
}

func (s LikeServer) ClearLike(ctx context.Context, req *like.MakeLikeRequest) (*like.MakeLikeResponse, error) {
	requestId := GetRequestId(ctx)
	requestString := MakeRequestString(ClearLikeMethod, req.PinId, req.UserId)
	s.logHandle(requestString, requestId)

	err := s.service.ClearLike(ctx, entity.PinID(req.PinId), entity.UserID(req.UserId))
	if err != emptyErrorInfo {
		s.logErr(requestString, err.GeneralErr, requestId)
		return &like.MakeLikeResponse{Valid: false, LocalError: int64(errs.ErrorCodes[err.LocalErr].LocalCode)}, nil
	}
	return &like.MakeLikeResponse{Valid: true}, nil
}

func MakeGetUsersLikedResponse(users []entity.UserResponse) *like.GetUsersLikedResponse {
	res := &like.GetUsersLikedResponse{Valid: true}
	for _, user := range users {
		res.Users = append(res.Users, &like.UserResponse{
			UserId:    int64(user.UserId),
			Email:     user.Email,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarURL,
		})
	}
	return res
}

func (s LikeServer) GetUsersLiked(ctx context.Context, req *like.GetUsersLikedRequest) (*like.GetUsersLikedResponse, error) {
	requestId := GetRequestId(ctx)
	requestString := MakeRequestString(GetUsersLikedMethod, req.PinId, req.Limit)
	s.logHandle(requestString, requestId)

	res, err := s.service.GetUsersLiked(ctx, entity.PinID(req.PinId), int(req.Limit))
	if err != emptyErrorInfo {
		s.logErr(requestString, err.GeneralErr, requestId)
		return &like.GetUsersLikedResponse{Valid: false, LocalError: int64(errs.ErrorCodes[err.LocalErr].LocalCode)}, nil
	}
	return MakeGetUsersLikedResponse(res.Users), nil
}

func (s LikeServer) CheckIsLiked(ctx context.Context, req *like.CheckIsLikedRequest) (*like.CheckIsLikedResponse, error) {
	requestId := GetRequestId(ctx)
	requestString := MakeRequestString(CheckIsLikedMethod, req.PinId, req.UserId)
	s.logHandle(requestString, requestId)

	res, err := s.service.CheckIsLiked(ctx, entity.PinID(req.PinId), entity.UserID(req.UserId))
	if err != nil {
		s.logErr(requestString, err, requestId)
		return &like.CheckIsLikedResponse{Valid: false, LocalError: 11}, nil
	}
	return &like.CheckIsLikedResponse{Valid: true, Liked: res}, nil
}

func MakeGetFavoritesResponseByUsers(pins entity.FeedPins) *like.GetFavoritesResponse {
	res := &like.GetFavoritesResponse{Valid: true}
	for _, pin := range pins.Pins {
		res.Pins = append(res.Pins, &like.FeedPin{
			PinId:      int64(pin.PinId),
			ContentUrl: pin.ContentUrl,
			Author: &like.PinAuthor{
				UserId:    int64(pin.PinAuthor.UserId),
				Nickname:  pin.PinAuthor.Nickname,
				AvatarUrl: pin.PinAuthor.AvatarURL,
			},
		})
	}
	return res
}

func (s LikeServer) GetFavorites(ctx context.Context, req *like.GetFavoritesRequest) (*like.GetFavoritesResponse, error) {
	requestId := GetRequestId(ctx)
	requestString := MakeRequestString(GetFavoritesMethod, req.UserId, req.Limit, req.Offset)
	s.logHandle(requestString, requestId)

	res, err := s.service.GetFavorites(ctx, entity.UserID(req.UserId), int(req.Limit), int(req.Offset))
	if err != emptyErrorInfo {
		s.logErr(requestString, err.GeneralErr, requestId)
		return &like.GetFavoritesResponse{Valid: false, LocalError: int64(errs.ErrorCodes[err.LocalErr].LocalCode)}, nil
	}
	return MakeGetFavoritesResponseByUsers(res), nil
}

func NewLikeServerForTests(s service.IService, l *zap.Logger) LikeServer {
	return LikeServer{service: s, logger: l}
}
