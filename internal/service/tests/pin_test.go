package test_service

import (
	"context"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_repository "harmonica/mocks/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockGetPinByIDArgs struct {
	Ctx   context.Context
	PinId entity.PinID
}
type mockGetPinByIDReturn struct {
	Pin entity.PinPageResponse
	Err error
}
type mockGetPinByID struct {
	Args   mockGetPinByIDArgs
	Return mockGetPinByIDReturn
}

var GetPinByIDCorrectValues = []mockGetPinByID{
	{
		Args: mockGetPinByIDArgs{
			Ctx:   context.Background(),
			PinId: 1,
		},
		Return: mockGetPinByIDReturn{
			Pin: entity.PinPageResponse{},
			Err: nil,
		},
	},
}
var GetPinByIDUncorrectValues = []mockGetPinByID{
	{
		Args: mockGetPinByIDArgs{
			Ctx:   context.Background(),
			PinId: 1,
		},
		Return: mockGetPinByIDReturn{
			Pin: entity.PinPageResponse{},
			Err: errs.ErrDBInternal,
		},
	},
}

func TestCorrectGetPinById(t *testing.T) {
	type funcArgs struct {
		Ctx   context.Context
		PinId entity.PinID
	}
	type funcReturn struct {
		Pin entity.PinPageResponse
		Err errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockGetPinByID     mockGetPinByID
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
		waitCheckCall      bool
		CheckReturn        bool
		ErrCheckReturn     error
	}
	tests := []test{
		{
			Name:           "Correct work test 1",
			MockGetPinByID: GetPinByIDCorrectValues[0],
			FuncArgs: funcArgs{
				Ctx:   context.Background(),
				PinId: 1,
			},
			ExpectedFuncReturn: funcReturn{},
			waitCheckCall:      true,
		},
		{
			Name:           "Error work test 1",
			MockGetPinByID: GetPinByIDUncorrectValues[0],
			FuncArgs: funcArgs{
				Ctx:   context.Background(),
				PinId: 1,
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
		{
			Name:           "Error work test 2",
			MockGetPinByID: GetPinByIDCorrectValues[0],
			FuncArgs: funcArgs{
				Ctx:   context.Background(),
				PinId: 1,
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			ErrCheckReturn: errs.ErrDBInternal,
			waitCheckCall:  true,
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().GetPinById(test.MockGetPinByID.Args.Ctx, test.MockGetPinByID.Args.PinId).Return(
			test.MockGetPinByID.Return.Pin, test.MockGetPinByID.Return.Err)
		if test.waitCheckCall {
			repo.EXPECT().CheckIsLiked(test.MockGetPinByID.Args.Ctx, test.MockGetPinByID.Args.PinId, entity.UserID(1)).Return(
				test.CheckReturn, test.ErrCheckReturn)
		}
		service := service.NewService(repo)
		res, err := service.GetPinById(test.FuncArgs.Ctx, test.FuncArgs.PinId, entity.UserID(1))
		assert.Equal(t, test.ExpectedFuncReturn.Pin, res)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestCreatePin(t *testing.T) {
	type mockArgs struct {
		Ctx context.Context
		Pin entity.Pin
	}
	type mockReturn struct {
		Pin entity.PinID
		Err error
	}
	type funcArgs struct {
		Ctx context.Context
		Pin entity.Pin
	}
	type funcReturn struct {
		Pin entity.PinPageResponse
		Err errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockGetPinByID     mockGetPinByID
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
		NotExpectGetPin    bool
	}

	tests := []test{
		{
			Name:           "Correct test 1",
			MockGetPinByID: GetPinByIDCorrectValues[0],
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{},
			},
			MockReturn: mockReturn{
				entity.PinID(1),
				nil,
			},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{},
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
			},
		},
		{
			Name:           "Uncorrect test 1",
			MockGetPinByID: GetPinByIDUncorrectValues[0],
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{},
			},
			MockReturn: mockReturn{
				entity.PinID(1),
				nil,
			},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{},
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
		{
			Name:           "Uncorrect test 2",
			MockGetPinByID: GetPinByIDCorrectValues[0],
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{},
			},
			MockReturn: mockReturn{
				entity.PinID(1),
				errs.ErrPermissionDenied,
			},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{},
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrPermissionDenied,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			NotExpectGetPin: true,
		},
	}

	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)

	for _, test := range tests {
		if !test.NotExpectGetPin {
			repo.EXPECT().GetPinById(test.MockGetPinByID.Args.Ctx, test.MockGetPinByID.Args.PinId).Return(
				test.MockGetPinByID.Return.Pin, test.MockGetPinByID.Return.Err)
		}
		repo.EXPECT().CreatePin(test.MockArgs.Ctx, test.MockArgs.Pin).Return(test.MockReturn.Pin, test.MockReturn.Err)
		service := service.NewService(repo)
		res, err := service.CreatePin(test.FuncArgs.Ctx, test.FuncArgs.Pin)
		assert.Equal(t, test.ExpectedFuncReturn.Pin, res)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestUpdatePin(t *testing.T) {
	type mockArgs struct {
		Ctx context.Context
		Pin entity.Pin
	}
	type mockReturn struct {
		Err error
	}
	type funcArgs struct {
		Ctx context.Context
		Pin entity.Pin
	}
	type funcReturn struct {
		Pin entity.PinPageResponse
		Err errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockGetPinByID     []mockGetPinByID
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
		NotExpectUpdatePin bool
	}
	tests := []test{
		{
			Name:           "Correct test 1",
			MockGetPinByID: []mockGetPinByID{GetPinByIDCorrectValues[0], GetPinByIDCorrectValues[0]},
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			MockReturn: mockReturn{
				nil,
			},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name:           "Uncorrect test 1",
			MockGetPinByID: []mockGetPinByID{GetPinByIDCorrectValues[0]},
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			MockReturn: mockReturn{
				errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			NotExpectUpdatePin: false,
		},
		{
			Name:           "Uncorrect test 2",
			MockGetPinByID: []mockGetPinByID{GetPinByIDUncorrectValues[0]},
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			NotExpectUpdatePin: true,
		},
		{
			Name:           "Uncorrect test 3",
			MockGetPinByID: []mockGetPinByID{GetPinByIDCorrectValues[0]},
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1), AuthorId: entity.UserID(10)},
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: nil,
					LocalErr:   errs.ErrPermissionDenied,
				},
			},
			NotExpectUpdatePin: true,
		},
		{
			Name:           "Uncorrect test 3",
			MockGetPinByID: []mockGetPinByID{GetPinByIDCorrectValues[0], GetPinByIDUncorrectValues[0]},
			MockArgs: mockArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			ExpectedFuncReturn: funcReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{
					GeneralErr: nil,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)

	for _, test := range tests {
		for _, request := range test.MockGetPinByID {
			repo.EXPECT().GetPinById(request.Args.Ctx, request.Args.PinId).Return(
				request.Return.Pin, request.Return.Err)
		}
		fmt.Println(test.Name)
		if !test.NotExpectUpdatePin {
			repo.EXPECT().UpdatePin(test.MockArgs.Ctx, test.MockArgs.Pin).Return(test.MockReturn.Err)
		}
		service := service.NewService(repo)
		res, err := service.UpdatePin(test.FuncArgs.Ctx, test.FuncArgs.Pin)
		assert.Equal(t, test.ExpectedFuncReturn.Pin, res)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestDeletePin(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		PinID entity.PinID
	}
	type mockReturn struct {
		Err error
	}
	type funcArgs struct {
		Ctx context.Context
		Pin entity.Pin
	}
	type funcReturn struct {
		Err errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockGetPinByID     mockGetPinByID
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
		NotExpectDeletePin bool
	}

	tests := []test{
		{
			Name:           "Correct test 1",
			MockGetPinByID: GetPinByIDCorrectValues[0],
			MockArgs: mockArgs{
				context.Background(),
				entity.PinID(1),
			},
			MockReturn: mockReturn{
				nil,
			},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name:           "Uncorrect test 1",
			MockGetPinByID: GetPinByIDCorrectValues[0],
			MockArgs: mockArgs{
				context.Background(),
				entity.PinID(1),
			},
			MockReturn: mockReturn{
				errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
		{
			Name:           "Uncorrect test 2",
			MockGetPinByID: GetPinByIDUncorrectValues[0],
			MockArgs:       mockArgs{},
			MockReturn:     mockReturn{},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1)},
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			NotExpectDeletePin: true,
		},
		{
			Name:           "Uncorrect test 3",
			MockGetPinByID: GetPinByIDCorrectValues[0],
			MockArgs: mockArgs{
				context.Background(),
				entity.PinID(1),
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				context.Background(),
				entity.Pin{PinId: entity.PinID(1), AuthorId: entity.UserID(100)},
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: nil,
					LocalErr:   errs.ErrPermissionDenied,
				},
			},
			NotExpectDeletePin: true,
		},
	}

	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)

	for _, test := range tests {
		repo.EXPECT().GetPinById(test.MockGetPinByID.Args.Ctx, test.MockGetPinByID.Args.PinId).Return(
			test.MockGetPinByID.Return.Pin, test.MockGetPinByID.Return.Err)
		if !test.NotExpectDeletePin {
			repo.EXPECT().DeletePin(test.MockArgs.Ctx, test.MockArgs.PinID).Return(test.MockReturn.Err)
		}
		service := service.NewService(repo)
		err := service.DeletePin(test.FuncArgs.Ctx, test.FuncArgs.Pin)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestGetUserPins(t *testing.T) {

	type mockArgs struct {
		Ctx      context.Context
		Nickname string
		UserID   entity.UserID
		Limit    int
		Offset   int
	}
	type mockReturn struct {
		UserResp entity.User
		PinsResp entity.UserPins
		UserErr  error
		PinsErr  error
	}
	type funcArgs struct {
		Ctx      context.Context
		Nickname string
		Limit    int
		Offset   int
	}
	type funcReturn struct {
		Pins entity.UserPins
		Err  errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
		NotExpectGetPins   bool
	}

	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "Nickname123",
				UserID:   entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				UserResp: entity.User{UserID: entity.UserID(1), Nickname: "Nickname123"},
				PinsResp: entity.UserPins{},
				UserErr:  nil,
				PinsErr:  nil,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "Nickname123",
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "Nickname123",
				UserID:   entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				UserResp: entity.User{UserID: entity.UserID(1), Nickname: "Nickname123"},
				PinsResp: entity.UserPins{},
				UserErr:  nil,
				PinsErr:  errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "Nickname123",
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
		{
			Name: "Uncorrect test 2",
			MockArgs: mockArgs{
				Ctx:      context.Background(),
				Nickname: "Nickname123",
				UserID:   entity.UserID(1),
				Limit:    10,
				Offset:   10,
			},
			MockReturn: mockReturn{
				UserResp: entity.User{UserID: entity.UserID(1), Nickname: "Nickname123"},
				PinsResp: entity.UserPins{},
				UserErr:  errs.ErrDBInternal,
				PinsErr:  nil,
			},
			FuncArgs: funcArgs{
				Ctx:      context.Background(),
				Nickname: "Nickname123",
				Limit:    10,
				Offset:   10,
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			NotExpectGetPins: true,
		},
	}

	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().GetUserByNickname(test.MockArgs.Ctx, test.MockArgs.Nickname).Return(
			test.MockReturn.UserResp, test.MockReturn.UserErr)
		if !test.NotExpectGetPins {
			repo.EXPECT().GetUserPins(test.MockArgs.Ctx, test.MockArgs.UserID, test.MockArgs.Limit, test.MockArgs.Offset).Return(
				test.MockReturn.PinsResp, test.MockReturn.PinsErr)
		}
		service := service.NewService(repo)
		res, err := service.GetUserPins(test.FuncArgs.Ctx, test.FuncArgs.Nickname, test.FuncArgs.Limit, test.FuncArgs.Offset)
		assert.Equal(t, test.ExpectedFuncReturn.Pins, res)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestGetFeedPins(t *testing.T) {
	type mockArgs struct {
		Ctx    context.Context
		Limit  int
		Offset int
	}
	type mockReturn struct {
		PinsResp entity.FeedPins
		Err      error
	}
	type funcArgs struct {
		Ctx    context.Context
		Limit  int
		Offset int
	}
	type funcReturn struct {
		Pins entity.FeedPins
		Err  errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
	}

	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				Limit:  10,
				Offset: 10,
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				Ctx:    context.Background(),
				Limit:  10,
				Offset: 10,
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				Limit:  10,
				Offset: 10,
			},
			MockReturn: mockReturn{
				PinsResp: entity.FeedPins{},
				Err:      errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:    context.Background(),
				Limit:  10,
				Offset: 10,
			},
			ExpectedFuncReturn: funcReturn{
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().GetFeedPins(test.MockArgs.Ctx, test.MockArgs.Limit, test.MockArgs.Offset).Return(
			test.MockReturn.PinsResp, test.MockReturn.Err)
		service := service.NewService(repo)
		res, err := service.GetFeedPins(test.FuncArgs.Ctx, test.FuncArgs.Limit, test.FuncArgs.Offset)
		assert.Equal(t, test.ExpectedFuncReturn.Pins, res)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}
