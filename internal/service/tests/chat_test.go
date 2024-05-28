package test_service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"
)

func TestService_CreateMessage(t *testing.T) {
	type Args struct {
		Message entity.Message
	}
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		Message entity.Message
		Draft   entity.Draft
	}
	type ExpectedMockReturn struct {
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().CreateMessage(ctx, mockArgs.Message).Return(mockReturn.Error)
		repo.EXPECT().UpdateDraft(ctx, mockArgs.Draft).Return(mockReturn.Error).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{Message: entity.Message{
				SenderId:   1,
				ReceiverId: 2,
				Text:       "Hello, World!",
			}},
			expectedReturn: ExpectedReturn{},
			expectedMockArgs: ExpectedMockArgs{
				Message: entity.Message{
					SenderId:   1,
					ReceiverId: 2,
					Text:       "Hello, World!",
				},
				Draft: entity.Draft{
					SenderId:   1,
					ReceiverId: 2,
				},
			},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name: "Error test case 1",
			args: Args{Message: entity.Message{
				SenderId:   1,
				ReceiverId: 999,
				Text:       "Test message",
			}},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: &pq.Error{Code: service.ForeignKeyViolationErrCode}, LocalErr: errs.ErrForeignKeyViolation},
			},
			expectedMockArgs: ExpectedMockArgs{Message: entity.Message{
				SenderId:   1,
				ReceiverId: 999,
				Text:       "Test message",
			}},
			expectedMockReturn: ExpectedMockReturn{
				Error: &pq.Error{Code: service.ForeignKeyViolationErrCode},
			},
		},
		{
			name: "Error test case 2",
			args: Args{Message: entity.Message{
				SenderId:   1,
				ReceiverId: 999,
				Text:       "Test message",
			}},
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{Message: entity.Message{
				SenderId:   1,
				ReceiverId: 999,
				Text:       "Test message",
			}},
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
			s := service.NewService(repo, likeClient)
			errInfo := s.CreateMessage(context.Background(), testCase.args.Message)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetMessages(t *testing.T) {
	type Args struct {
		UserId1 entity.UserID
		UserId2 entity.UserID
	}
	type ExpectedReturn struct {
		Messages  entity.Messages
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		UserId1 entity.UserID
		UserId2 entity.UserID
	}
	type ExpectedMockReturn struct {
		Messages entity.Messages
		Draft    entity.DraftResponse
		Error1   error
		Error2   error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetMessages(ctx, mockArgs.UserId1, mockArgs.UserId2).Return(mockReturn.Messages, mockReturn.Error1)
		repo.EXPECT().GetDraft(ctx, mockArgs.UserId1, mockArgs.UserId2).Return(mockReturn.Draft, mockReturn.Error2).MaxTimes(1)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{UserId1: 1, UserId2: 2},
			expectedReturn: ExpectedReturn{
				Messages: entity.Messages{Messages: []entity.MessageResponse{}},
			},
			expectedMockArgs: ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{
				Messages: entity.Messages{Messages: []entity.MessageResponse{}},
			},
		},
		{
			name: "Error test case 1",
			args: Args{UserId1: 1, UserId2: 2},
			expectedReturn: ExpectedReturn{
				Messages:  entity.Messages{},
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{
				Messages: entity.Messages{},
				Error1:   errs.ErrDBInternal,
			},
		},
		{
			name: "Error test case 2",
			args: Args{UserId1: 1, UserId2: 2},
			expectedReturn: ExpectedReturn{
				Messages:  entity.Messages{},
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{UserId1: 1, UserId2: 2},
			expectedMockReturn: ExpectedMockReturn{
				Messages: entity.Messages{},
				Error2:   errs.ErrDBInternal,
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
			s := service.NewService(repo, likeClient)
			messages, errInfo := s.GetMessages(context.Background(), testCase.args.UserId1, testCase.args.UserId2)
			assert.Equal(t, testCase.expectedReturn.Messages, messages)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_GetUserChats(t *testing.T) {
	type Args struct {
		UserId entity.UserID
	}
	type ExpectedReturn struct {
		UserChats entity.UserChats
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockArgs struct {
		UserId entity.UserID
	}
	type ExpectedMockReturn struct {
		UserChats entity.UserChats
		Error     error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUserChats(ctx, mockArgs.UserId).Return(mockReturn.UserChats, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: Args{UserId: 1},
			expectedReturn: ExpectedReturn{
				UserChats: entity.UserChats{Chats: []entity.UserChat{}},
			},
			expectedMockArgs: ExpectedMockArgs{UserId: 1},
			expectedMockReturn: ExpectedMockReturn{
				UserChats: entity.UserChats{Chats: []entity.UserChat{}},
			},
		},
		{
			name: "Error test case 1",
			args: Args{UserId: 1},
			expectedReturn: ExpectedReturn{
				UserChats: entity.UserChats{},
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockArgs: ExpectedMockArgs{UserId: 1},
			expectedMockReturn: ExpectedMockReturn{
				UserChats: entity.UserChats{},
				Error:     errs.ErrDBInternal,
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
			s := service.NewService(repo, likeClient)
			userChats, errInfo := s.GetUserChats(context.Background(), testCase.args.UserId)
			assert.Equal(t, testCase.expectedReturn.UserChats, userChats)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
