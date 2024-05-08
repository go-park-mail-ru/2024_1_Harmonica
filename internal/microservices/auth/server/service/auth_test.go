package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	mock_repository "harmonica/mocks/repository"
	"testing"
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
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := NewService(repo)
			user, errInfo := s.GetUserByEmail(context.Background(), testCase.args)
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
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := NewService(repo)
			user, errInfo := s.GetUserById(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.User, user)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
