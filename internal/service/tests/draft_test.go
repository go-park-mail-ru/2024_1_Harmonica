package test_service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"
)

func TestService_UpdateDraft(t *testing.T) {
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Error error
	}
	draft := entity.Draft{SenderId: 2, ReceiverId: 123, Text: "Test draft"}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockReturn ExpectedMockReturn) {
		repo.EXPECT().UpdateDraft(ctx, draft).Return(mockReturn.Error)
	}
	testTable := []struct {
		name               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:               "OK test case 1",
			expectedReturn:     ExpectedReturn{},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name: "Error test case 1",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
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
			mockBehaviour(repo, context.Background(), testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			errInfo := s.UpdateDraft(context.Background(), draft)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
