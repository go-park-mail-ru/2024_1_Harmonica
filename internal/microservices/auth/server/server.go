package server

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	auth "harmonica/internal/microservices/auth/proto"
	"harmonica/internal/microservices/auth/server/service"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

var (
	Sessions            sync.Map
	sessionTTL          = 24 * time.Hour
	sessionsCleanupTime = 6 * time.Hour
	emptyUser           = entity.User{}
	emptyErrorInfo      = errs.ErrorInfo{}
)

type AuthorizationServer struct {
	service service.IService
	auth.AuthorizationServer
}

func NewAuthorizationServer(s *service.Service) AuthorizationServer {
	return AuthorizationServer{service: s}
}

func (s AuthorizationServer) CheckSession(ctx context.Context, req *auth.CheckSessionRequest) (*auth.CheckSessionResponse, error) {
	ses, exists := Sessions.Load(req.Session)
	if !exists || ses.(Session).IsExpired() {
		if exists {
			Sessions.Delete(req.Session)
		}
		return &auth.CheckSessionResponse{Valid: false, LocalError: 2}, nil
	}
	userId := ses.(Session).UserId
	return &auth.CheckSessionResponse{Valid: true, UserId: int64(userId)}, nil
}

func (s AuthorizationServer) Login(ctx context.Context, req *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	if len(metadata.ValueFromIncomingContext(ctx, "request_id")) == 0 {
		return &auth.LoginUserResponse{Valid: false, LocalError: 7}, nil
	}
	ctx = context.WithValue(ctx, "request_id", metadata.ValueFromIncomingContext(ctx, "request_id")[0])
	if !ValidateEmail(req.Email) || !ValidatePassword(req.Password) {
		return &auth.LoginUserResponse{LocalError: 5, Valid: false}, nil
	}
	loggedInUser, errInfo := s.service.GetUserByEmail(ctx, req.Email) // 1service

	if errInfo != emptyErrorInfo {
		return &auth.LoginUserResponse{LocalError: int64(errs.ErrorCodes[errInfo.LocalErr].LocalCode), Valid: false}, nil
	}
	if loggedInUser == emptyUser {
		return &auth.LoginUserResponse{LocalError: 7, Valid: false}, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(loggedInUser.Password), []byte(req.Password))
	if err != nil {
		return &auth.LoginUserResponse{LocalError: 8, Valid: false}, nil
	}
	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(sessionTTL)
	ses := Session{
		UserId: loggedInUser.UserID,
		Expiry: expiresAt,
	}
	Sessions.Store(newSessionToken, ses)
	return GetLoginResponseUserByUser(loggedInUser, newSessionToken, expiresAt.Format(time.RFC3339Nano)), nil
}

func GetLoginResponseUserByUser(user entity.User, token string, expiresAt string) *auth.LoginUserResponse {
	return &auth.LoginUserResponse{
		Valid:           true,
		UserId:          int64(user.UserID),
		Email:           user.Email,
		Nickname:        user.Nickname,
		Password:        user.Password,
		AvatarURL:       user.AvatarURL,
		RegisterAt:      user.RegisterAt.Format(time.RFC3339Nano),
		NewSessionToken: token,
		ExpiresAt:       expiresAt,
	}
}

func (s AuthorizationServer) IsAuth(ctx context.Context, req *auth.Empty) (*auth.IsAuthResponse, error) {
	if len(metadata.ValueFromIncomingContext(ctx, "user_id")) == 0 ||
		len(metadata.ValueFromIncomingContext(ctx, "request_id")) == 0 {
		return MakeLocalErrorIsAuth(2), nil
	}

	ctx = context.WithValue(ctx, "user_id", metadata.ValueFromIncomingContext(ctx, "user_id")[0])
	ctx = context.WithValue(ctx, "request_id", metadata.ValueFromIncomingContext(ctx, "request_id")[0])

	userIdFromSession := metadata.ValueFromIncomingContext(ctx, "user_id")[0]
	if userIdFromSession == "" {
		return MakeLocalErrorIsAuth(2), nil
	}
	id, _ := strconv.Atoi(userIdFromSession)
	userId := entity.UserID(id)
	user, errInfo := s.service.GetUserById(ctx, userId) // 2service
	if errInfo != emptyErrorInfo {
		return MakeLocalErrorIsAuth(int64(errs.ErrorCodes[errInfo.LocalErr].LocalCode)), nil
	}

	if user == emptyUser {
		return MakeLocalErrorIsAuth(2), nil
	}

	return &auth.IsAuthResponse{User: GetUserResponseByUser(user), IsAuthorized: true, Valid: true}, nil
}

func MakeLocalErrorIsAuth(localCode int64) *auth.IsAuthResponse {
	return &auth.IsAuthResponse{
		IsAuthorized: false,
		LocalError:   localCode,
		Valid:        false,
	}
}

func GetUserResponseByUser(user entity.User) *auth.IsAuthUserResponse {
	return &auth.IsAuthUserResponse{UserId: int64(user.UserID), Email: user.Email, Nickname: user.Nickname, AvatarURL: user.AvatarURL}
}

func (s AuthorizationServer) Logout(ctx context.Context, in *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	_, exists := Sessions.Load(in.SessionToken)
	if !exists {
		return nil, nil
	}
	Sessions.Delete(in.SessionToken)
	return nil, nil
}
