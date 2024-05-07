package test_service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"
)

func TestService_AddSubscriptionToUser(t *testing.T) {
	type Args struct {
		UserId1 entity.UserID
		UserId2 entity.UserID
	}
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		UserId1 entity.UserID
		UserId2 entity.UserID
	}
	type ExpectedMockReturn struct {
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().AddSubscriptionToUser(ctx, mockArgs.UserId1, mockArgs.UserId2).Return(mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:               "OK test case 1",
			args:               Args{UserId1: 1, UserId2: 2},
			expectedReturn:     ExpectedReturn{},
			expectedMockArgs:   ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name: "Error test case 1",
			args: Args{UserId1: 1, UserId2: 2},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: Args{UserId1: 1, UserId2: 2},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: &pq.Error{Code: service.UniqueViolationErrCode}, LocalErr: errs.ErrDBUniqueViolation},
			},
			expectedMockArgs: ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{
				Error: &pq.Error{Code: service.UniqueViolationErrCode},
			},
		},
		{
			name: "Error test case 3",
			args: Args{UserId1: 1, UserId2: 2},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: &pq.Error{Code: service.ForeignKeyViolationErrCode}, LocalErr: errs.ErrForeignKeyViolation},
			},
			expectedMockArgs: ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{
				Error: &pq.Error{Code: service.ForeignKeyViolationErrCode},
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
			errInfo := s.AddSubscriptionToUser(context.Background(), testCase.args.UserId1, testCase.args.UserId2)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_DeleteSubscriptionToUser(t *testing.T) {
	type Args struct {
		UserId1 entity.UserID
		UserId2 entity.UserID
	}
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		UserId1 entity.UserID
		UserId2 entity.UserID
	}
	type ExpectedMockReturn struct {
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().DeleteSubscriptionToUser(ctx, mockArgs.UserId1, mockArgs.UserId2).Return(mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:               "OK test case 1",
			args:               Args{UserId1: 1, UserId2: 2},
			expectedReturn:     ExpectedReturn{},
			expectedMockArgs:   ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name: "Error test case 1",
			args: Args{UserId1: 1, UserId2: 2},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{UserId1: 1, UserId2: 2},
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
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			errInfo := s.DeleteSubscriptionToUser(context.Background(), testCase.args.UserId1, testCase.args.UserId2)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetUserSubscribers(t *testing.T) {
	type Args struct {
		UserId entity.UserID
	}
	type ExpectedReturn struct {
		Subscribers entity.UserSubscribers
		ErrorInfo   errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		UserId entity.UserID
	}
	type ExpectedMockReturn struct {
		Subscribers entity.UserSubscribers
		Error       error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserSubscribers(ctx, mockArgs.UserId).Return(mockReturn.Subscribers, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{UserId: 1},
			expectedReturn: ExpectedReturn{
				Subscribers: entity.UserSubscribers{Subscribers: []entity.UserResponse{{UserId: 2, Email: "example@mail.com", Nickname: "nickname"}}},
				ErrorInfo:   errs.ErrorInfo{},
			},
			expectedMockArgs: ExpectedMockArgs{UserId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Subscribers: entity.UserSubscribers{Subscribers: []entity.UserResponse{{UserId: 2, Email: "example@mail.com", Nickname: "nickname"}}},
				Error:       nil,
			},
		},
		{
			name: "Error test case 1",
			args: Args{UserId: 1},
			expectedReturn: ExpectedReturn{
				Subscribers: entity.UserSubscribers{},
				ErrorInfo:   errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{UserId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Subscribers: entity.UserSubscribers{},
				Error:       errs.ErrDBInternal,
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
			subscribers, errInfo := s.GetUserSubscribers(context.Background(), testCase.args.UserId)
			assert.Equal(t, testCase.expectedReturn.Subscribers, subscribers)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetUserSubscriptions(t *testing.T) {
	type Args struct {
		UserId entity.UserID
	}
	type ExpectedReturn struct {
		Subscriptions entity.UserSubscriptions
		ErrorInfo     errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		UserId entity.UserID
	}
	type ExpectedMockReturn struct {
		Subscriptions entity.UserSubscriptions
		Error         error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserSubscriptions(ctx, mockArgs.UserId).Return(mockReturn.Subscriptions, mockReturn.Error)
	}

	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{UserId: 1},
			expectedReturn: ExpectedReturn{
				Subscriptions: entity.UserSubscriptions{
					Subscriptions: []entity.UserResponse{
						{UserId: 2, Email: "example@mail.com", Nickname: "nickname", AvatarURL: "url", AvatarDX: 100, AvatarDY: 100},
					},
				},
				ErrorInfo: errs.ErrorInfo{},
			},
			expectedMockArgs: ExpectedMockArgs{UserId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Subscriptions: entity.UserSubscriptions{
					Subscriptions: []entity.UserResponse{
						{UserId: 2, Email: "example@mail.com", Nickname: "nickname", AvatarURL: "url", AvatarDX: 100, AvatarDY: 100},
					},
				},
				Error: nil,
			},
		},
		{
			name: "Error test case 1",
			args: Args{UserId: 1},
			expectedReturn: ExpectedReturn{
				Subscriptions: entity.UserSubscriptions{},
				ErrorInfo:     errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{UserId: 1},
			expectedMockReturn: ExpectedMockReturn{
				Subscriptions: entity.UserSubscriptions{},
				Error:         errs.ErrDBInternal,
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
			subscriptions, errInfo := s.GetUserSubscriptions(context.Background(), testCase.args.UserId)
			assert.Equal(t, testCase.expectedReturn.Subscriptions, subscriptions)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
