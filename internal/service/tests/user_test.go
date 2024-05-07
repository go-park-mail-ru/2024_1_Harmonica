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
		repo.EXPECT().GetUserByEmail(ctx, mockArgs).Return(
			mockReturn.User, mockReturn.Error)
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
			service := service.NewService(repo, likeClient)
			user, errInfo := service.GetUserByEmail(context.Background(), testCase.args)
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
		repo.EXPECT().GetUserByNickname(ctx, mockArgs).Return(
			mockReturn.User, mockReturn.Error)
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
			service := service.NewService(repo, likeClient)
			user, errInfo := service.GetUserByNickname(context.Background(), testCase.args)
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
			service := service.NewService(repo, likeClient)
			user, errInfo := service.GetUserById(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.User, user)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_RegisterUser(t *testing.T) {
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
		repo.EXPECT().GetUserByEmail(ctx, mockArgs).Return(
			mockReturn.User, mockReturn.Error).AnyTimes()
		repo.EXPECT().GetUserByNickname(ctx, mockArgs).Return(
			mockReturn.User, mockReturn.Error).AnyTimes()
		repo.EXPECT().RegisterUser(ctx, mockArgs)
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
			service := service.NewService(repo, likeClient)
			user, errInfo := service.GetUserById(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.User, user)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

/*

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
