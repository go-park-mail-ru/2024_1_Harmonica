package test_service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUserByEmail(t *testing.T) {
	type ExpectedReturn struct {
		User      entity.User
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		User  entity.User
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs string, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserByEmail(ctx, mockArgs).Return(mockReturn.User, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               string
		expectedReturn     ExpectedReturn
		expectedMockArgs   string
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: "email@mail.ru",
			expectedReturn: ExpectedReturn{
				User: entity.User{UserID: 1, Email: "email@mail.ru"},
			},
			expectedMockArgs: "email@mail.ru",
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{UserID: 1, Email: "email@mail.ru"},
			},
		},
		{
			name: "Error test case 1",
			args: "email@mail.ru",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: "email@mail.ru",
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
			user, errInfo := s.GetUserByEmail(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.User, user)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetUserByNickname(t *testing.T) {
	type ExpectedReturn struct {
		User      entity.User
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		User  entity.User
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs string, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserByNickname(ctx, mockArgs).Return(mockReturn.User, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               string
		expectedReturn     ExpectedReturn
		expectedMockArgs   string
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: "nickname1",
			expectedReturn: ExpectedReturn{
				User: entity.User{UserID: 1, Nickname: "nickname1"},
			},
			expectedMockArgs: "nickname1",
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{UserID: 1, Nickname: "nickname1"},
			},
		},
		{
			name: "Error test case 1",
			args: "nickname1",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: "nickname1",
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
			user, errInfo := s.GetUserByNickname(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.User, user)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetUserById(t *testing.T) {
	type ExpectedReturn struct {
		User      entity.User
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		User  entity.User
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.UserID, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserById(ctx, mockArgs).Return(
			mockReturn.User, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               entity.UserID
		expectedReturn     ExpectedReturn
		expectedMockArgs   entity.UserID
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: 1,
			expectedReturn: ExpectedReturn{
				User: entity.User{UserID: 1},
			},
			expectedMockArgs: 1,
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{UserID: 1},
			},
		},
		{
			name: "Error test case 1",
			args: 1,
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: 1,
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
			user, errInfo := s.GetUserById(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.User, user)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_RegisterUser(t *testing.T) {
	type ExpectedMockReturn struct {
		User   entity.User
		Error1 error
		Error2 error
		Error3 error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.User, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserByEmail(ctx, mockArgs.Email).Return(mockReturn.User, mockReturn.Error1)
		repo.EXPECT().GetUserByNickname(ctx, mockArgs.Nickname).Return(mockReturn.User, mockReturn.Error2).MaxTimes(1)
		repo.EXPECT().RegisterUser(ctx, mockArgs).Return(mockReturn.Error3).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               entity.User
		expectedReturn     []errs.ErrorInfo
		expectedMockArgs   entity.User
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:             "OK test case 1",
			args:             entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
		},
		{
			name: "Error test case 1",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: []errs.ErrorInfo{
				{LocalErr: errs.ErrDBUniqueEmail},
				{LocalErr: errs.ErrDBUniqueNickname},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			},
		},
		{
			name: "Error test case 2",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: []errs.ErrorInfo{
				{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 3",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: []errs.ErrorInfo{
				{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				Error2: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 4",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: []errs.ErrorInfo{
				{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				Error3: errs.ErrDBInternal,
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
			errsInfo := s.RegisterUser(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn, errsInfo)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	type ExpectedReturn struct {
		User      entity.User
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		User   entity.User
		Error1 error
		Error2 error
		Error3 error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.User, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserByNickname(ctx, mockArgs.Nickname).Return(mockReturn.User, mockReturn.Error1)
		repo.EXPECT().UpdateUser(ctx, mockArgs).Return(mockReturn.Error2).MaxTimes(1)
		repo.EXPECT().GetUserById(ctx, mockArgs.UserID).Return(mockReturn.User, mockReturn.Error3).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               entity.User
		expectedReturn     ExpectedReturn
		expectedMockArgs   entity.User
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				User: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			},
		},
		{
			name: "Error test case 1",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrDBUniqueNickname},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{UserID: 123, Email: "email@mail.ru", Nickname: "nickname"},
			},
		},
		{
			name: "Error test case 2",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 3",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				Error2: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 4",
			args: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: entity.User{UserID: 1, Email: "email@mail.ru", Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				Error3: errs.ErrDBInternal,
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
			user, errInfo := s.UpdateUser(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.User, user)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetUserProfileByNickname(t *testing.T) {
	type ExpectedReturn struct {
		UserProfile entity.UserProfileResponse
		ErrorInfo   errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		User        entity.User
		UserProfile entity.UserProfileResponse
		Error1      error
		Error2      error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs entity.User, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserByNickname(ctx, mockArgs.Nickname).Return(mockReturn.User, mockReturn.Error1)
		repo.EXPECT().GetSubscriptionsInfo(ctx, mockArgs.UserID, mockArgs.UserID).Return(mockReturn.UserProfile, mockReturn.Error2).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               entity.User
		expectedReturn     ExpectedReturn
		expectedMockArgs   entity.User
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: entity.User{UserID: 1, Nickname: "nickname"},
			expectedReturn: ExpectedReturn{
				UserProfile: entity.UserProfileResponse{
					User: entity.UserResponse{
						UserId:   1,
						Email:    "email@mail.ru",
						Nickname: "nickname",
					},
					SubscriptionsCount: 5,
					SubscribersCount:   10,
					IsOwner:            true,
				},
			},
			expectedMockArgs: entity.User{UserID: 1, Nickname: "nickname"},
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{
					UserID:   1,
					Email:    "email@mail.ru",
					Nickname: "nickname",
				},
				UserProfile: entity.UserProfileResponse{
					User: entity.UserResponse{
						UserId:   1,
						Email:    "email@mail.ru",
						Nickname: "nickname",
					},
					SubscriptionsCount: 5,
					SubscribersCount:   10,
					IsOwner:            true,
				},
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
				User: entity.User{
					UserID:   1,
					Email:    "email@mail.ru",
					Nickname: "nickname",
				},
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
			userProfile, errInfo := s.GetUserProfileByNickname(context.Background(), testCase.args.Nickname, testCase.args.UserID)
			assert.Equal(t, testCase.expectedReturn.UserProfile, userProfile)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
