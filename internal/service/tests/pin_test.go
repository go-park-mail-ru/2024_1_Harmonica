package test_service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	like "harmonica/internal/microservices/like/proto"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"
)

const (
	Limit  = 10
	Offset = 10
)

func TestService_GetFeedPins(t *testing.T) {
	type ExpectedReturn struct {
		FeedPins  entity.FeedPins
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		FeedPins entity.FeedPins
		Error    error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		limit int, offset int, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetFeedPins(ctx, limit, offset).Return(mockReturn.FeedPins, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			expectedReturn: ExpectedReturn{
				FeedPins: entity.FeedPins{
					Pins: []entity.FeedPinResponse{
						{
							PinId:      1,
							ContentUrl: "bobrdobr.ru/1",
							PinAuthor: entity.PinAuthor{
								UserId:    123,
								Nickname:  "valera",
								AvatarURL: "bobrdobr.users.ru/123",
							},
						},
					},
				},
			},
			expectedMockReturn: ExpectedMockReturn{
				FeedPins: entity.FeedPins{
					Pins: []entity.FeedPinResponse{
						{
							PinId:      1,
							ContentUrl: "bobrdobr.ru/1",
							PinAuthor: entity.PinAuthor{
								UserId:    123,
								Nickname:  "valera",
								AvatarURL: "bobrdobr.users.ru/123",
							},
						},
					},
				},
			},
		},
		{
			name: "Error test case 1",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), Limit, Offset, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			pins, errInfo := s.GetFeedPins(context.Background(), Limit, Offset)
			assert.Equal(t, testCase.expectedReturn.FeedPins, pins)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetSubscriptionsFeedPins(t *testing.T) {
	type ExpectedReturn struct {
		FeedPins  entity.FeedPins
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		FeedPins entity.FeedPins
		Error    error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		userId entity.UserID, limit int, offset int, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetSubscriptionsFeedPins(ctx, userId, limit, offset).Return(mockReturn.FeedPins, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               entity.UserID
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: 1,
		},
		{
			name: "Error test case 1",
			args: 1,
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.args, Limit, Offset, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			pins, errInfo := s.GetSubscriptionsFeedPins(context.Background(), testCase.args, Limit, Offset)
			assert.Equal(t, testCase.expectedReturn.FeedPins, pins)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetUserPins(t *testing.T) {
	type ExpectedReturn struct {
		UserPins  entity.UserPins
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		User     entity.User
		UserPins entity.UserPins
		Error1   error
		Error2   error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.User, limit int, offset int, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserByNickname(ctx, mockArgs.Nickname).Return(mockReturn.User, mockReturn.Error1)
		repo.EXPECT().GetUserPins(ctx, mockArgs.UserID, limit, offset).Return(mockReturn.UserPins, mockReturn.Error2).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               entity.User
		expectedReturn     ExpectedReturn
		expectedMockArgs   entity.User
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:             "OK test case 1",
			args:             entity.User{UserID: 1, Nickname: "nickname"},
			expectedReturn:   ExpectedReturn{},
			expectedMockArgs: entity.User{UserID: 1, Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{UserID: 1, Nickname: "nickname"},
			},
		},
		{
			name: "Error test case 1",
			args: entity.User{UserID: 1, Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: entity.User{UserID: 1, Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				User:   entity.User{UserID: 1, Nickname: "nickname"},
				Error2: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, Limit, Offset, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			pins, errInfo := s.GetUserPins(context.Background(), testCase.args.Nickname, Limit, Offset)
			assert.Equal(t, testCase.expectedReturn.UserPins, pins)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetPinById(t *testing.T) {
	type Args struct {
		PinId  entity.PinID
		UserId entity.UserID
	}
	type ExpectedReturn struct {
		PinPageResponse entity.PinPageResponse
		ErrorInfo       errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Pin             entity.PinPageResponse
		CheckIsLikedRes *like.CheckIsLikedResponse
		Error1          error
		Error2          error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, likeService *mock_proto.MockLikeClient, ctx context.Context, mockArgs Args, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetPinById(ctx, mockArgs.PinId).Return(mockReturn.Pin, mockReturn.Error1)
		likeService.EXPECT().CheckIsLiked(ctx, &like.CheckIsLikedRequest{PinId: int64(mockArgs.PinId),
			UserId: int64(mockArgs.UserId)}).Return(mockReturn.CheckIsLikedRes, mockReturn.Error2).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{PinId: 1, UserId: 1},
			expectedReturn: ExpectedReturn{
				PinPageResponse: entity.PinPageResponse{
					PinId:   1,
					Title:   "Test Pin",
					IsLiked: true,
				},
			},
			expectedMockReturn: ExpectedMockReturn{
				Pin:             entity.PinPageResponse{PinId: 1, Title: "Test Pin"},
				CheckIsLikedRes: &like.CheckIsLikedResponse{Valid: true, Liked: true},
			},
		},
		{
			name: "Error test case 1",
			args: Args{PinId: 2, UserId: 1},
			expectedReturn: ExpectedReturn{
				PinPageResponse: entity.PinPageResponse{},
				ErrorInfo:       errs.ErrorInfo{GeneralErr: errors.New("pin not found"), LocalErr: errs.ErrElementNotExist},
			},
			expectedMockReturn: ExpectedMockReturn{
				Pin:             entity.PinPageResponse{},
				CheckIsLikedRes: &like.CheckIsLikedResponse{Valid: true, Liked: true},
				Error1:          errors.New("pin not found"),
			},
		},
		{
			name: "Error test case 2",
			args: Args{PinId: 1, UserId: 1},
			expectedReturn: ExpectedReturn{
				PinPageResponse: entity.PinPageResponse{},
				ErrorInfo:       errs.ErrorInfo{LocalErr: errs.ErrGRPCWentWrong},
			},
			expectedMockReturn: ExpectedMockReturn{
				Pin:             entity.PinPageResponse{},
				CheckIsLikedRes: &like.CheckIsLikedResponse{Valid: true, Liked: true},
				Error2:          errors.New("grpc error"),
			},
		},
		{
			name: "Error test case 3",
			args: Args{PinId: 1, UserId: 1},
			expectedReturn: ExpectedReturn{
				PinPageResponse: entity.PinPageResponse{},
				ErrorInfo:       errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				CheckIsLikedRes: &like.CheckIsLikedResponse{Valid: false, LocalError: int64(errs.ErrorCodes[errs.ErrDBInternal].LocalCode)},
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeService := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, likeService, context.Background(), testCase.args, testCase.expectedMockReturn)
			s := service.NewService(repo, likeService)
			pin, errInfo := s.GetPinById(context.Background(), testCase.args.PinId, testCase.args.UserId)
			assert.Equal(t, testCase.expectedReturn.PinPageResponse, pin)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_CreatePin(t *testing.T) {
	type ExpectedReturn struct {
		Pin       entity.PinPageResponse
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Pin    entity.PinPageResponse
		Error1 error
		Error2 error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.Pin, mockReturn ExpectedMockReturn) {
		repo.EXPECT().CreatePin(ctx, mockArgs).Return(mockReturn.Pin.PinId, mockReturn.Error1)
		repo.EXPECT().GetPinById(ctx, mockArgs.PinId).Return(mockReturn.Pin, mockReturn.Error2).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               entity.Pin
		expectedReturn     ExpectedReturn
		expectedMockArgs   entity.Pin
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				Pin: entity.PinPageResponse{PinId: 1},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin: entity.PinPageResponse{PinId: 1},
			},
		},
		{
			name: "Error test case 1",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin:    entity.PinPageResponse{PinId: 1},
				Error2: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			pin, errInfo := s.CreatePin(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.Pin, pin)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_UpdatePin(t *testing.T) {
	type ExpectedReturn struct {
		Pin       entity.PinPageResponse
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Pin    entity.PinPageResponse
		Error1 error
		Error2 error
		Error3 error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.Pin, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetPinById(ctx, mockArgs.PinId).Return(mockReturn.Pin, mockReturn.Error1).MaxTimes(1)
		repo.EXPECT().UpdatePin(ctx, mockArgs).Return(mockReturn.Error2).MaxTimes(1)
		repo.EXPECT().GetPinById(ctx, mockArgs.PinId).Return(mockReturn.Pin, mockReturn.Error3).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               entity.Pin
		expectedReturn     ExpectedReturn
		expectedMockArgs   entity.Pin
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				Pin: entity.PinPageResponse{PinId: 1},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin: entity.PinPageResponse{PinId: 1},
			},
		},
		{
			name: "Error test case 1",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrElementNotExist},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin:    entity.PinPageResponse{PinId: 1},
				Error2: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 3",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin:    entity.PinPageResponse{PinId: 1},
				Error3: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 4",
			args: entity.Pin{PinId: 1, AuthorId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrPermissionDenied},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin: entity.PinPageResponse{PinId: 1, PinAuthor: entity.PinAuthor{UserId: 123}},
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			pin, errInfo := s.UpdatePin(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.Pin, pin)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_DeletePin(t *testing.T) {
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Pin    entity.PinPageResponse
		Error1 error
		Error2 error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.Pin, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetPinById(ctx, mockArgs.PinId).Return(mockReturn.Pin, mockReturn.Error1)
		repo.EXPECT().DeletePin(ctx, mockArgs.PinId).Return(mockReturn.Error2).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               entity.Pin
		expectedReturn     ExpectedReturn
		expectedMockArgs   entity.Pin
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:             "OK test case 1",
			args:             entity.Pin{PinId: 1},
			expectedReturn:   ExpectedReturn{},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin: entity.PinPageResponse{PinId: 1},
			},
		},
		{
			name: "Error test case 1",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrElementNotExist},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: entity.Pin{PinId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin:    entity.PinPageResponse{PinId: 1},
				Error2: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 3",
			args: entity.Pin{PinId: 1, AuthorId: 1},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrPermissionDenied},
			},
			expectedMockArgs: entity.Pin{PinId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Pin: entity.PinPageResponse{PinId: 1, PinAuthor: entity.PinAuthor{UserId: 123}},
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			likeClient := mock_proto.NewMockLikeClient(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			errInfo := s.DeletePin(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
