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

func TestGetUserByEmail(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		Email string
	}
	type mockReturn struct {
		User entity.User
		Err  error
	}
	type funcArgs struct {
		Ctx   context.Context
		Email string
	}
	type funcReturn struct {
		User entity.User
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
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: "email@dot.com",
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				Ctx:   context.Background(),
				Email: "email@dot.com",
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: "email@dot.com",
			},
			MockReturn: mockReturn{
				Err: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:   context.Background(),
				Email: "email@dot.com",
			},
			ExpectedFuncReturn: funcReturn{
				entity.User{},
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
		repo.EXPECT().GetUserByEmail(test.MockArgs.Ctx, test.MockArgs.Email).Return(
			test.MockReturn.User, test.MockReturn.Err)
		service := service.NewService(repo)
		user, err := service.GetUserByEmail(test.FuncArgs.Ctx, test.FuncArgs.Email)
		assert.Equal(t, test.ExpectedFuncReturn.User, user)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestGetUserByNickname(t *testing.T) {
	type mockArgs struct {
		Ctx  context.Context
		Nick string
	}
	type mockReturn struct {
		User entity.User
		Err  error
	}
	type funcArgs struct {
		Ctx  context.Context
		Nick string
	}
	type funcReturn struct {
		User entity.User
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
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:  context.Background(),
				Nick: "nickdot.com",
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				Ctx:  context.Background(),
				Nick: "nickdot.com",
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:  context.Background(),
				Nick: "nickdot.com",
			},
			MockReturn: mockReturn{
				Err: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx:  context.Background(),
				Nick: "nickdot.com",
			},
			ExpectedFuncReturn: funcReturn{
				entity.User{},
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
		repo.EXPECT().GetUserByNickname(test.MockArgs.Ctx, test.MockArgs.Nick).Return(
			test.MockReturn.User, test.MockReturn.Err)
		service := service.NewService(repo)
		user, err := service.GetUserByNickname(test.FuncArgs.Ctx, test.FuncArgs.Nick)
		assert.Equal(t, test.ExpectedFuncReturn.User, user)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}

func TestGetUserById(t *testing.T) {
	type mockArgs struct {
		Ctx context.Context
		Id  entity.UserID
	}
	type mockReturn struct {
		User entity.User
		Err  error
	}
	type funcArgs struct {
		Ctx context.Context
		Id  entity.UserID
	}
	type funcReturn struct {
		User entity.User
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
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx: context.Background(),
				Id:  entity.UserID(1),
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				Id:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{},
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx: context.Background(),
				Id:  entity.UserID(1),
			},
			MockReturn: mockReturn{
				Err: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				Id:  entity.UserID(1),
			},
			ExpectedFuncReturn: funcReturn{
				entity.User{},
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
		repo.EXPECT().GetUserById(test.MockArgs.Ctx, test.MockArgs.Id).Return(
			test.MockReturn.User, test.MockReturn.Err)
		service := service.NewService(repo)
		user, err := service.GetUserById(test.FuncArgs.Ctx, test.FuncArgs.Id)
		assert.Equal(t, test.ExpectedFuncReturn.User, user)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}
func TestRegisterUser(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		Email string
		Nick  string
		User  entity.User
	}
	type mockReturn struct {
		User   entity.User
		GetErr error
		RegErr error
	}
	type funcArgs struct {
		Ctx  context.Context
		User entity.User
	}
	type funcReturn struct {
		errs []errs.ErrorInfo
	}
	type test struct {
		Name               string
		MockArgs           mockArgs
		MockReturn         mockReturn
		FuncArgs           funcArgs
		ExpectedFuncReturn funcReturn
		WaiRegCall         bool
		WaitNickCall       bool
	}
	tests := []test{
		{
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: "email.dot@com",
				Nick:  "nickoNick",
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			MockReturn: mockReturn{},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			ExpectedFuncReturn: funcReturn{},
			WaiRegCall:         true,
			WaitNickCall:       true,
		},
		{
			Name: "Uncorrect work test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: "email.dot@com",
				Nick:  "nickoNick",
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			MockReturn: mockReturn{
				GetErr: errs.ErrDBInternal,
				RegErr: nil,
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			ExpectedFuncReturn: funcReturn{
				[]errs.ErrorInfo{
					{
						GeneralErr: errs.ErrDBInternal,
						LocalErr:   errs.ErrDBInternal,
					},
				},
			},
			WaiRegCall:   false,
			WaitNickCall: true,
		},
		{
			Name: "Uncorrect work test 2",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: "email.dot@com",
				Nick:  "nickoNick",
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			MockReturn: mockReturn{
				RegErr: errs.ErrDBInternal,
				User:   entity.User{},
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			ExpectedFuncReturn: funcReturn{
				[]errs.ErrorInfo{
					{
						GeneralErr: errs.ErrDBInternal,
						LocalErr:   errs.ErrDBInternal,
					},
				},
			},
			WaiRegCall:   false,
			WaitNickCall: false,
		},
		{
			Name: "Uncorrect work test 3",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: "email.dot@com",
				Nick:  "nickoNick",
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			MockReturn: mockReturn{
				RegErr: errs.ErrDBInternal,
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			ExpectedFuncReturn: funcReturn{
				[]errs.ErrorInfo{
					{
						GeneralErr: nil,
						LocalErr:   errs.ErrDBUniqueEmail,
					},
					{
						GeneralErr: nil,
						LocalErr:   errs.ErrDBUniqueNickname,
					},
				},
			},
			WaiRegCall:   false,
			WaitNickCall: true,
		},
		{
			Name: "Uncorrect work test 4",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: "email.dot@com",
				Nick:  "nickoNick",
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			MockReturn: mockReturn{
				RegErr: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Email:    "email.dot@com",
					Nickname: "nickoNick",
				},
			},
			ExpectedFuncReturn: funcReturn{
				[]errs.ErrorInfo{
					{
						GeneralErr: errs.ErrDBInternal,
						LocalErr:   errs.ErrDBInternal,
					},
				},
			},
			WaiRegCall:   true,
			WaitNickCall: true,
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		repo.EXPECT().GetUserByEmail(test.MockArgs.Ctx, test.MockArgs.Email).Return(
			test.MockReturn.User, test.MockReturn.GetErr)
		if test.WaitNickCall {
			repo.EXPECT().GetUserByNickname(test.MockArgs.Ctx, test.MockArgs.Nick).Return(
				test.MockReturn.User, test.MockReturn.GetErr)
		}
		if test.WaiRegCall {
			repo.EXPECT().RegisterUser(test.MockArgs.Ctx, test.MockArgs.User).Return(test.MockReturn.RegErr)
		}
		service := service.NewService(repo)
		errs := service.RegisterUser(test.FuncArgs.Ctx, test.FuncArgs.User)
		assert.Equal(t, test.ExpectedFuncReturn.errs, errs)
	}
}

func TestUpdateUser(t *testing.T) {
	type mockArgs struct {
		Ctx        context.Context
		UpdateUser entity.User
	}
	type mockReturn struct {
		User             entity.User
		ErrGetById       error
		ErrGetByNickanme error
		ErrUpdating      error
	}
	type funcArgs struct {
		Ctx  context.Context
		User entity.User
	}
	type funcReturn struct {
		User entity.User
		Err  errs.ErrorInfo
	}
	type test struct {
		Name                string
		MockArgs            mockArgs
		MockReturn          mockReturn
		FuncArgs            funcArgs
		ExpectedFuncReturn  funcReturn
		ExcpectIDCall       bool
		ExcpectNicknameCall bool
		ExcpectUpdateCall   bool
	}
	tests := []test{
		{
			Name: "Correct work test 1",
			MockArgs: mockArgs{
				Ctx: context.Background(),
				UpdateUser: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			MockReturn: mockReturn{
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			ExpectedFuncReturn: funcReturn{
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			ExcpectIDCall:       true,
			ExcpectNicknameCall: true,
			ExcpectUpdateCall:   true,
		},
		{
			Name: "Uncorrect work test 2",
			MockArgs: mockArgs{
				Ctx: context.Background(),
				UpdateUser: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			MockReturn: mockReturn{
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
				ErrGetByNickanme: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			ExpectedFuncReturn: funcReturn{
				User: entity.User{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			ExcpectIDCall:       false,
			ExcpectNicknameCall: true,
			ExcpectUpdateCall:   false,
		},
		{
			Name: "Uncorrect work test 2",
			MockArgs: mockArgs{
				Ctx: context.Background(),
				UpdateUser: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			MockReturn: mockReturn{
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(2),
				},
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			ExpectedFuncReturn: funcReturn{
				User: entity.User{},
				Err: errs.ErrorInfo{
					GeneralErr: nil,
					LocalErr:   errs.ErrDBUniqueNickname,
				},
			},
			ExcpectIDCall:       false,
			ExcpectNicknameCall: true,
			ExcpectUpdateCall:   false,
		},
		{
			Name: "Uncorrect work test 3",
			MockArgs: mockArgs{
				Ctx: context.Background(),
				UpdateUser: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			MockReturn: mockReturn{
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
				ErrUpdating: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			ExpectedFuncReturn: funcReturn{
				User: entity.User{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			ExcpectIDCall:       false,
			ExcpectNicknameCall: true,
			ExcpectUpdateCall:   true,
		},
		{
			Name: "Uncorrect work test 4",
			MockArgs: mockArgs{
				Ctx: context.Background(),
				UpdateUser: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			MockReturn: mockReturn{
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
				ErrGetById: errs.ErrDBInternal,
			},
			FuncArgs: funcArgs{
				Ctx: context.Background(),
				User: entity.User{
					Nickname: "nickNNName11",
					UserID:   entity.UserID(1),
				},
			},
			ExpectedFuncReturn: funcReturn{
				User: entity.User{},
				Err: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			ExcpectIDCall:       true,
			ExcpectNicknameCall: true,
			ExcpectUpdateCall:   true,
		},
	}
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockIRepository(ctrl)
	for _, test := range tests {
		if test.ExcpectIDCall {
			returnUser := test.MockReturn.User
			if test.MockReturn.ErrGetById != nil {
				returnUser = entity.User{}
			}
			repo.EXPECT().GetUserById(test.MockArgs.Ctx, test.MockArgs.UpdateUser.UserID).Return(
				returnUser, test.MockReturn.ErrGetById)
		}
		if test.ExcpectNicknameCall {
			returnUser := test.MockReturn.User
			if test.MockReturn.ErrGetByNickanme != nil {
				returnUser = entity.User{}
			}
			repo.EXPECT().GetUserByNickname(test.MockArgs.Ctx, test.MockArgs.UpdateUser.Nickname).Return(
				returnUser, test.MockReturn.ErrGetByNickanme)
		}
		if test.ExcpectUpdateCall {
			repo.EXPECT().UpdateUser(test.MockArgs.Ctx, test.MockArgs.UpdateUser).Return(
				test.MockReturn.ErrUpdating)
		}
		service := service.NewService(repo)
		user, err := service.UpdateUser(test.FuncArgs.Ctx, test.FuncArgs.User)
		assert.Equal(t, test.ExpectedFuncReturn.User, user)
		assert.Equal(t, test.ExpectedFuncReturn.Err, err)
	}
}
*/
