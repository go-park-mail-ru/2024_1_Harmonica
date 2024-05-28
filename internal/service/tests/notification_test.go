package test_service

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/service"
	mock_proto "harmonica/mocks/microservices/like/proto"
	mock_repository "harmonica/mocks/repository"
	"testing"
)

func TestService_GetUnreadNotifications(t *testing.T) {
	type ExpectedReturn struct {
		Notifications entity.Notifications
		ErrorInfo     errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Notifications entity.Notifications
		Error         error
	}
	userId := entity.UserID(2)
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetUnreadNotifications(ctx, userId).Return(mockReturn.Notifications, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			expectedReturn: ExpectedReturn{
				Notifications: entity.Notifications{
					Notifications: []entity.NotificationResponse{
						{
							UserId: userId,
							Type:   "new_pin",
						},
					},
				},
			},
			expectedMockReturn: ExpectedMockReturn{
				Notifications: entity.Notifications{
					Notifications: []entity.NotificationResponse{
						{
							UserId: userId,
							Type:   "new_pin",
						},
					},
				},
			},
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
			notifications, errInfo := s.GetUnreadNotifications(context.Background(), userId)
			assert.Equal(t, testCase.expectedReturn.Notifications, notifications)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_CreateNotification(t *testing.T) {
	type ExpectedReturn struct {
		NotificationId entity.NotificationID
		ErrorInfo      errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		NotificationId entity.NotificationID
		Error          error
	}
	notification := entity.Notification{NotificationId: 22}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockReturn ExpectedMockReturn) {
		repo.EXPECT().CreateNotification(ctx, notification).Return(mockReturn.NotificationId, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:               "OK test case 1",
			expectedReturn:     ExpectedReturn{NotificationId: notification.NotificationId},
			expectedMockReturn: ExpectedMockReturn{NotificationId: notification.NotificationId},
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
			notificationId, errInfo := s.CreateNotification(context.Background(), notification)
			assert.Equal(t, testCase.expectedReturn.NotificationId, notificationId)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_ReadNotification(t *testing.T) {
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type MockArgs struct {
		NotificationId entity.NotificationID
	}
	type ExpectedMockReturn struct {
		Notification entity.Notification
		Error1       error
		Error2       error
	}
	userId := entity.UserID(2)
	notification := entity.Notification{NotificationId: 22}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetNotificationById(ctx, notification.NotificationId).Return(mockReturn.Notification, mockReturn.Error1)
		repo.EXPECT().ReadNotification(ctx, notification).Return(mockReturn.Error2).AnyTimes()
	}
	testTable := []struct {
		name               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
		mockArgs           MockArgs
	}{
		{
			name:               "OK test case 1",
			expectedReturn:     ExpectedReturn{},
			expectedMockReturn: ExpectedMockReturn{Notification: notification},
			mockArgs:           MockArgs{NotificationId: notification.NotificationId},
		},
		{
			name:               "OK test case 2",
			expectedReturn:     ExpectedReturn{},
			expectedMockReturn: ExpectedMockReturn{Error1: sql.ErrNoRows},
			mockArgs:           MockArgs{NotificationId: notification.NotificationId},
		},
		{
			name: "Error test case 1",
			expectedReturn: ExpectedReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			expectedMockReturn: ExpectedMockReturn{Error1: errs.ErrDBInternal},
			mockArgs:           MockArgs{NotificationId: notification.NotificationId},
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
			errInfo := s.ReadNotification(context.Background(), notification.NotificationId, userId)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}

func TestService_ReadAllNotifications(t *testing.T) {
	type ExpectedReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type ExpectedMockReturn struct {
		Error error
	}
	userId := entity.UserID(2)
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context, mockReturn ExpectedMockReturn) {
		repo.EXPECT().ReadAllNotifications(ctx, userId).Return(mockReturn.Error)
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
			expectedMockReturn: ExpectedMockReturn{Error: errs.ErrDBInternal},
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
			errInfo := s.ReadAllNotifications(context.Background(), userId)
			assert.Equal(t, testCase.expectedReturn.ErrorInfo, errInfo)
		})
	}
}
