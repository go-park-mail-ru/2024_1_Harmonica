package test_service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"
)

func TestService_Search(t *testing.T) {
	type Args struct {
		Query string
	}
	type ExpectedReturn struct {
		SearchResult entity.SearchResult
		ErrorInfo    errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Users  []entity.SearchUser
		Pins   []entity.SearchPin
		Boards []entity.SearchBoard
		Error1 error
		Error2 error
		Error3 error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, query string, mockReturn ExpectedMockReturn) {
		repo.EXPECT().SearchForUsers(ctx, fmt.Sprintf(`%s%s%s`, "%", query, "%")).Return(mockReturn.Users, mockReturn.Error1)
		repo.EXPECT().SearchForPins(ctx, query).Return(mockReturn.Pins, mockReturn.Error2).MaxTimes(1)
		repo.EXPECT().SearchForBoards(ctx, query).Return(mockReturn.Boards, mockReturn.Error3).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{Query: "test"},
			expectedReturn: ExpectedReturn{
				SearchResult: entity.SearchResult{
					Users:  []entity.SearchUser{{UserId: 1, Nickname: "TestUser1"}},
					Pins:   []entity.SearchPin{{PinId: 2}},
					Boards: []entity.SearchBoard{{BoardId: 3}},
				},
			},
			expectedMockReturn: ExpectedMockReturn{
				Users:  []entity.SearchUser{{UserId: 1, Nickname: "TestUser1"}},
				Pins:   []entity.SearchPin{{PinId: 2}},
				Boards: []entity.SearchBoard{{BoardId: 3}},
			},
		},
		{
			name: "Error test case 1",
			args: Args{Query: "test"},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: Args{Query: "test"},
			expectedReturn: ExpectedReturn{
				SearchResult: entity.SearchResult{},
				ErrorInfo:    errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error2: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 3",
			args: Args{Query: "test"},
			expectedReturn: ExpectedReturn{
				SearchResult: entity.SearchResult{},
				ErrorInfo:    errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: errs.ErrDBInternal},
			},
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
			mockBehaviour(repo, context.Background(), testCase.args.Query, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			result, errInfo := s.Search(context.Background(), testCase.args.Query)
			assert.Equal(t, testCase.expectedReturn.SearchResult, result)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
