package server

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	like "harmonica/internal/microservices/like/proto"
	mock_service "harmonica/mocks/microservices/like/server/service"
	"testing"
)

func TestSetLike(t *testing.T) {
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		service.EXPECT().SetLike(ctx, mockArgs.PinID, mockArgs.UserID).Return(mockReturn.ErrorInfo)
	}
	testTable := []testStruct{
		{
			name: "Good test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				Response: &like.MakeLikeResponse{Valid: true},
			},
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
				Response: &like.MakeLikeResponse{Valid: false},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errors.New("some error")},
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			service, server := configureMock(t)
			mockBehaviour(service, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			response, err := server.SetLike(context.Background(), &like.MakeLikeRequest{
				PinId:  int64(testCase.args.PinID),
				UserId: int64(testCase.args.UserID),
			})
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

func TestClearLike(t *testing.T) {
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		service.EXPECT().ClearLike(ctx, mockArgs.PinID, mockArgs.UserID).Return(mockReturn.ErrorInfo)
	}
	testTable := configureTests()
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			service, server := configureMock(t)
			mockBehaviour(service, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			response, err := server.ClearLike(context.Background(), &like.MakeLikeRequest{
				PinId:  int64(testCase.args.PinID),
				UserId: int64(testCase.args.UserID),
			})
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

func TestGetUsersLiked(t *testing.T) {
	type Args struct {
		PinID entity.PinID
		Limit int64
	}
	type ExpectedReturn struct {
		Response *like.GetUsersLikedResponse
		Error    error
	}
	type ExpectedMockArgs struct {
		PinID entity.PinID
		Limit int64
	}
	type ExpectedMockReturn struct {
		Response  entity.UserList
		ErrorInfo errs.ErrorInfo
	}
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		service.EXPECT().GetUsersLiked(ctx, mockArgs.PinID, int(mockArgs.Limit)).Return(mockReturn.Response, mockReturn.ErrorInfo)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "Good test",
			args: Args{
				PinID: 1,
				Limit: 10,
			},
			expectedReturn: ExpectedReturn{
				Response: MakeGetUsersLikedResponse([]entity.UserResponse{{UserId: 1}}),
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID: 1,
				Limit: 10,
			},
			expectedMockReturn: ExpectedMockReturn{
				Response: entity.UserList{Users: []entity.UserResponse{{UserId: 1}}},
			},
		},
		{
			name: "Error test",
			args: Args{
				PinID: 1,
				Limit: 10,
			},
			expectedReturn: ExpectedReturn{
				Response: &like.GetUsersLikedResponse{Valid: false},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID: 1,
				Limit: 10,
			},
			expectedMockReturn: ExpectedMockReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errors.New("some error")},
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mock_service.NewMockIService(ctrl)
			server := NewLikeServerForTests(service, zap.NewNop())
			mockBehaviour(service, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			response, err := server.GetUsersLiked(context.Background(), &like.GetUsersLikedRequest{
				PinId: int64(testCase.args.PinID),
				Limit: testCase.args.Limit,
			})
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

func TestCheckIsLiked(t *testing.T) {
	type Args struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedReturn struct {
		Response *like.CheckIsLikedResponse
		Error    error
	}
	type ExpectedMockArgs struct {
		PinID  entity.PinID
		UserID entity.UserID
	}
	type ExpectedMockReturn struct {
		Response bool
		Error    error
	}
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		service.EXPECT().CheckIsLiked(ctx, mockArgs.PinID, mockArgs.UserID).Return(mockReturn.Response, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "Good test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				Response: &like.CheckIsLikedResponse{Valid: true, Liked: true},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				Response: true,
			},
		},
		{
			name: "Error test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				Response: &like.CheckIsLikedResponse{Valid: false, LocalError: 11},
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
			service := mock_service.NewMockIService(ctrl)
			server := NewLikeServerForTests(service, zap.NewNop())
			mockBehaviour(service, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			response, err := server.CheckIsLiked(context.Background(), &like.CheckIsLikedRequest{
				PinId:  int64(testCase.args.PinID),
				UserId: int64(testCase.args.UserID),
			})
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

func TestGetFavorites(t *testing.T) {
	type Args struct {
		UserID entity.UserID
		Limit  int64
		Offset int64
	}
	type ExpectedReturn struct {
		Response *like.GetFavoritesResponse
		Error    error
	}
	type ExpectedMockArgs struct {
		UserID entity.UserID
		Limit  int64
		Offset int64
	}
	type ExpectedMockReturn struct {
		Response  entity.FeedPins
		ErrorInfo errs.ErrorInfo
	}
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		service.EXPECT().GetFavorites(ctx, mockArgs.UserID, int(mockArgs.Limit), int(mockArgs.Offset)).
			Return(mockReturn.Response, mockReturn.ErrorInfo)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "Good test",
			args: Args{
				UserID: 1,
				Limit:  10,
				Offset: 0,
			},
			expectedReturn: ExpectedReturn{
				Response: &like.GetFavoritesResponse{Valid: true},
			},
			expectedMockArgs: ExpectedMockArgs{
				UserID: 1,
				Limit:  10,
				Offset: 0,
			},
			expectedMockReturn: ExpectedMockReturn{
				Response: entity.FeedPins{},
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
				Response: &like.GetFavoritesResponse{Valid: false},
			},
			expectedMockArgs: ExpectedMockArgs{
				UserID: 1,
				Limit:  10,
				Offset: 0,
			},
			expectedMockReturn: ExpectedMockReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errors.New("some error")},
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mock_service.NewMockIService(ctrl)
			likeServer := NewLikeServerForTests(service, zap.NewNop())
			mockBehaviour(service, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			response, err := likeServer.GetFavorites(context.Background(), &like.GetFavoritesRequest{
				UserId: int64(testCase.args.UserID),
				Limit:  testCase.args.Limit,
				Offset: testCase.args.Offset,
			})
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

// зачем это? - линтер.............
type Args struct {
	PinID  entity.PinID
	UserID entity.UserID
}
type ExpectedReturn struct {
	Response *like.MakeLikeResponse
	Error    error
}
type ExpectedMockArgs struct {
	PinID  entity.PinID
	UserID entity.UserID
}
type ExpectedMockReturn struct {
	ErrorInfo errs.ErrorInfo
}
type testStruct struct {
	name               string
	args               Args
	expectedReturn     ExpectedReturn
	expectedMockArgs   ExpectedMockArgs
	expectedMockReturn ExpectedMockReturn
}

func configureTests() []testStruct {
	tests := []testStruct{
		{
			name: "Good test",
			args: Args{
				PinID:  1,
				UserID: 1,
			},
			expectedReturn: ExpectedReturn{
				Response: &like.MakeLikeResponse{Valid: true},
			},
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
				Response: &like.MakeLikeResponse{Valid: false},
			},
			expectedMockArgs: ExpectedMockArgs{
				PinID:  1,
				UserID: 1,
			},
			expectedMockReturn: ExpectedMockReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errors.New("some error")},
			},
		},
	}
	return tests
}

func configureMock(t *testing.T) (*mock_service.MockIService, LikeServer) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock_service.NewMockIService(ctrl)
	server := NewLikeServerForTests(service, zap.NewNop())
	return service, server
}
