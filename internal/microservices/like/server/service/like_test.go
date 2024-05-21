package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	mock_repository "harmonica/mocks/microservices/like/server/repository"
	"testing"
)

func TestSetLike(t *testing.T) {
	type Args struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedMockReturn struct {
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().SetLike(ctx, mockArgs.PinID, mockArgs.UserID).Return(mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name: "Error test 1",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errors.New("some error"), LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errors.New("some error"),
			},
		},
		{
			name: "Error test 2",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: &pq.Error{Code: "23503"},
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			service := NewRepositoryService(repo)
			errInfo := service.SetLike(context.Background(), testCase.args.PinID, testCase.args.UserID)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestClearLike(t *testing.T) {
	type Args struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedMockReturn struct {
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().ClearLike(ctx, mockArgs.PinID, mockArgs.UserID).Return(mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name: "Error test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errors.New("some error"), LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errors.New("some error"),
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			service := NewRepositoryService(repo)
			errInfo := service.ClearLike(context.Background(), testCase.args.PinID, testCase.args.UserID)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestGetUsersLiked(t *testing.T) {
	type Args struct {
		PinID entity.PinID
		Limit int
	}
	type ExpectedReturn struct {
		UserList entity.UserList
		Error    errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		PinID entity.PinID
		Limit int
	}
	type ExpectedMockReturn struct {
		UserList entity.UserList
		Error    error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUsersLiked(ctx, mockArgs.PinID, mockArgs.Limit).Return(mockReturn.UserList, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test",
			args: Args{
				PinID: 1,
				Limit: 10,
			},
			expectedReturn: ExpectedReturn{
				UserList: entity.UserList{},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID: 1,
				Limit: 10,
			},
			expectedMockReturn: ExpectedMockReturn{
				UserList: entity.UserList{},
			},
		},
		{
			name: "Error test",
			args: Args{
				PinID: 1,
				Limit: 10,
			},
			expectedReturn: ExpectedReturn{
				Error: errs.ErrorInfo{GeneralErr: errors.New("some error"), LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID: 1,
				Limit: 10,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errors.New("some error"),
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			service := NewRepositoryService(repo)
			userList, errInfo := service.GetUsersLiked(context.Background(), testCase.args.PinID, testCase.args.Limit)
			assert.Equal(t, testCase.expectedReturn.UserList, userList)
			assert.Equal(t, testCase.expectedReturn.Error, errInfo)
		})
	}
}

func TestGetFavorites(t *testing.T) {
	type Args struct {
		UserID entity.UserID
		Limit  int
		Offset int
	}
	type ExpectedReturn struct {
		FeedPins entity.FeedPins
		Error    errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		UserID entity.UserID
		Limit  int
		Offset int
	}
	type ExpectedMockReturn struct {
		FeedPins entity.FeedPins
		Error    error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetFavorites(ctx, mockArgs.UserID, mockArgs.Limit, mockArgs.Offset).Return(mockReturn.FeedPins, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test",
			args: Args{
				UserID: 1,
				Limit:  10,
				Offset: 0,
			},
			expectedReturn: ExpectedReturn{
				FeedPins: entity.FeedPins{},
			},
			expectedMockArgs: ExpectedMockArgs{
				UserID: 1,
				Limit:  10,
				Offset: 0,
			},
			expectedMockReturn: ExpectedMockReturn{
				FeedPins: entity.FeedPins{},
			},
		},
		{
			name: "Error test",
			args: Args{
				UserID: 1,
				Limit:  10,
				Offset: 0,
			},
			expectedReturn: ExpectedReturn{
				Error: errs.ErrorInfo{GeneralErr: errors.New("some error"), LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{
				UserID: 1,
				Limit:  10,
				Offset: 0,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errors.New("some error"),
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			service := NewRepositoryService(repo)
			feedPins, errInfo := service.GetFavorites(context.Background(), testCase.args.UserID, testCase.args.Limit, testCase.args.Offset)
			assert.Equal(t, testCase.expectedReturn.FeedPins, feedPins)
			assert.Equal(t, testCase.expectedReturn.Error, errInfo)
		})
	}
}

func TestCheckIsLiked(t *testing.T) {
	type Args struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedReturn struct {
		IsLiked bool
		Error   error
	}
	type ExpectedMockArgs struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedMockReturn struct {
		IsLiked bool
		Error   error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().CheckIsLiked(ctx, mockArgs.PinID, mockArgs.UserID).Return(mockReturn.IsLiked, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				IsLiked: true,
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				IsLiked: true,
			},
		},
		{
			name: "Error test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				Error: errors.New("some error"),
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errors.New("some error"),
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			service := NewRepositoryService(repo)
			isLiked, err := service.CheckIsLiked(context.Background(), testCase.args.PinID, testCase.args.UserID)
			assert.Equal(t, testCase.expectedReturn.IsLiked, isLiked)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}
