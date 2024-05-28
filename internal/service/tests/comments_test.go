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

func TestService_AddComment(t *testing.T) {
	type ExpectedReturn struct {
		Pin       entity.PinPageResponse
		CommentId entity.CommentID
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Exists    bool
		CommentId entity.CommentID
		Pin       entity.PinPageResponse
		Error1    error
		Error2    error
		Error3    error
	}
	pinId := entity.PinID(1)
	pin := entity.PinPageResponse{PinId: pinId}
	userId := entity.UserID(2)
	commentId := entity.CommentID(2)
	myComment := "my comment"
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, comment string, mockReturn ExpectedMockReturn) {
		repo.EXPECT().CheckPinExistence(ctx, pinId).Return(mockReturn.Exists, mockReturn.Error1)
		repo.EXPECT().AddComment(ctx, comment, pinId, userId).Return(mockReturn.CommentId, mockReturn.Error2).AnyTimes()
		repo.EXPECT().GetPinById(ctx, pinId).Return(mockReturn.Pin, mockReturn.Error3).AnyTimes()
	}
	testTable := []struct {
		name               string
		comment            string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:    "OK test case 1",
			comment: myComment,
			expectedReturn: ExpectedReturn{
				Pin:       pin,
				CommentId: commentId,
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists:    true,
				CommentId: commentId,
				Pin:       pin,
			},
		},
		{
			name: "Error test case 1",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrNotFound},
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists: false,
			},
		},
		{
			name:    "Error test case 3",
			comment: "",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrEmptyComment},
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists: true,
			},
		},
		{
			name:    "Error test case 4",
			comment: myComment,
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists: true,
				Error2: errs.ErrDBInternal,
			},
		},
		{
			name:    "Error test case 5",
			comment: myComment,
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists: true,
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
			mockBehaviour(repo, context.Background(), testCase.comment, testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			curPin, curCommentId, errInfo := s.AddComment(context.Background(), testCase.comment, pinId, userId)
			assert.Equal(t, testCase.expectedReturn.Pin, curPin)
			assert.Equal(t, testCase.expectedReturn.CommentId, curCommentId)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetComments(t *testing.T) {
	type ExpectedReturn struct {
		Comments  entity.GetCommentsResponse
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Exists   bool
		Comments entity.GetCommentsResponse
		Error1   error
		Error2   error
	}
	pinId := entity.PinID(1)
	comments := entity.GetCommentsResponse{Comments: []entity.CommentResponse{{CommentId: 2}}}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockReturn ExpectedMockReturn) {
		repo.EXPECT().CheckPinExistence(ctx, pinId).Return(mockReturn.Exists, mockReturn.Error1)
		repo.EXPECT().GetComments(ctx, pinId).Return(mockReturn.Comments, mockReturn.Error2).AnyTimes()
	}
	testTable := []struct {
		name               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			expectedReturn: ExpectedReturn{
				Comments: comments,
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists:   true,
				Comments: comments,
			},
		},
		{
			name: "Error test case 1",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error1: errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{LocalErr: errs.ErrNotFound},
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists: false,
			},
		},
		{
			name: "Error test case 3",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{
				Exists: true,
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
			mockBehaviour(repo, context.Background(), testCase.expectedMockReturn)
			s := service.NewService(repo, likeClient)
			curComments, errInfo := s.GetComments(context.Background(), pinId)
			assert.Equal(t, testCase.expectedReturn.Comments, curComments)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
