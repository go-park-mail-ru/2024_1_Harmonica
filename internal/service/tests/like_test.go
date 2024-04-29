package test_service

/*
import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_repository "harmonica/mocks/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetLike(t *testing.T) {
	type mockArgs struct {
		Ctx    context.Context
		PinId  entity.PinID
		UserId entity.UserID
	}
	type mockReturn struct {
		Err error
	}
	type funcArgs struct {
		Ctx    context.Context
		PinId  entity.PinID
		UserId entity.UserID
	}
	type funcReturn struct {
		Err errs.ErrorInfo
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
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			MockReturn: mockReturn{
				Err: nil,
			},
			FuncArgs: funcArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			MockReturn: mockReturn{
				Err: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().SetLike(test.MockArgs.Ctx, test.MockArgs.PinId, test.MockArgs.UserId).Return(
			test.MockReturn.Err)
		service := service.NewService(repo)
		err := service.SetLike(test.FuncArgs.Ctx, test.FuncArgs.PinId, test.FuncArgs.UserId)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestClearLike(t *testing.T) {
	type mockArgs struct {
		Ctx    context.Context
		PinId  entity.PinID
		UserId entity.UserID
	}
	type mockReturn struct {
		Err error
	}
	type funcArgs struct {
		Ctx    context.Context
		PinId  entity.PinID
		UserId entity.UserID
	}
	type funcReturn struct {
		Err errs.ErrorInfo
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
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			MockReturn: mockReturn{
				Err: nil,
			},
			FuncArgs: funcArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			MockReturn: mockReturn{
				Err: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().ClearLike(test.MockArgs.Ctx, test.MockArgs.PinId, test.MockArgs.UserId).Return(
			test.MockReturn.Err)
		service := service.NewService(repo)
		err := service.ClearLike(test.FuncArgs.Ctx, test.FuncArgs.PinId, test.FuncArgs.UserId)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestGetUsersLiked(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		PinId entity.PinID
		Limit int
	}
	type mockReturn struct {
		ListUsers entity.UserList
		Err       error
	}
	type funcArgs struct {
		Ctx   context.Context
		PinId entity.PinID
		Limit int
	}
	type funcReturn struct {
		ListUsers entity.UserList
		Err       errs.ErrorInfo
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
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				PinId: entity.PinID(1),
				Limit: 10,
			},
			MockReturn: mockReturn{
				Err: nil,
			},
			FuncArgs: funcArgs{
				Ctx:   context.Background(),
				PinId: entity.PinID(1),
				Limit: 10,
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				PinId: entity.PinID(1),
				Limit: 10,
			},
			MockReturn: mockReturn{
				Err: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:   context.Background(),
				PinId: entity.PinID(1),
				Limit: 10,
			},
			ExpectedFuncReturn: funcReturn{
				entity.UserList{},
				errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().GetUsersLiked(test.MockArgs.Ctx, test.MockArgs.PinId, test.MockArgs.Limit).Return(
			entity.UserList{}, test.MockReturn.Err)
		service := service.NewService(repo)
		list, err := service.GetUsersLiked(test.FuncArgs.Ctx, test.FuncArgs.PinId, test.FuncArgs.Limit)
		assert.Equal(t, test.ExpectedFuncReturn.ListUsers, list)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}
*/
