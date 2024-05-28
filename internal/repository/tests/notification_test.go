package test_repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/repository"
	"testing"
	"time"
)

func TestDBRepository_CreateNotification(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name           string
		notification   entity.Notification
		setupMocks     func()
		expectedResult entity.NotificationID
		expectedErr    error
	}{
		{
			name: "Create new pin notification",
			notification: entity.Notification{
				Type:              entity.NotificationTypeNewPin,
				TriggeredByUserId: 1,
				PinId:             2,
			},
			setupMocks: func() {
				mock.ExpectExec(repository.QueryCreatePinNotification).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedResult: 0,
			expectedErr:    nil,
		},
		{
			name: "Create base notification",
			notification: entity.Notification{
				UserId:            1,
				Type:              "comment",
				TriggeredByUserId: 2,
				PinId:             3,
				CommentId:         4,
				MessageId:         5,
			},
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCreateBaseNotification).
					WithArgs(1, "comment", 2, 3, 4, 5).
					WillReturnRows(sqlmock.NewRows([]string{"notification_id"}).AddRow(1))
			},
			expectedResult: 1,
			expectedErr:    nil,
		},
		{
			name: "Error creating notification",
			notification: entity.Notification{
				UserId:            1,
				Type:              "comment",
				TriggeredByUserId: 2,
				PinId:             3,
				CommentId:         4,
				MessageId:         5,
			},
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCreateBaseNotification).
					WithArgs(1, "comment", 2, 3, 4, 5).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedResult: 0,
			expectedErr:    errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			result, err := repo.CreateNotification(CtxWithRequestId, tc.notification)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestDBRepository_GetNotificationById(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()
	testTime := time.Now()

	tests := []struct {
		name           string
		notificationId entity.NotificationID
		setupMocks     func()
		expectedResult entity.Notification
		expectedErr    error
	}{
		{
			name:           "OK test",
			notificationId: 1,
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetNotificationById).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"notification_id", "user_id", "type", "triggered_by_user_id",
						"pin_id", "comment_id", "message_id", "created_at"}).
						AddRow(1, 1, "comment", 2, 3, 4, 5, testTime))
			},
			expectedResult: entity.Notification{
				NotificationId:    1,
				UserId:            1,
				Type:              "comment",
				TriggeredByUserId: 2,
				PinId:             3,
				CommentId:         4,
				MessageId:         5,
				CreatedAt:         testTime,
			},
			expectedErr: nil,
		},
		{
			name:           "Error test",
			notificationId: 1,
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetNotificationById).
					WithArgs(1).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedResult: entity.Notification{},
			expectedErr:    errs.ErrDBInternal,
		},
		{
			name:           "No rows found test",
			notificationId: 1,
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetNotificationById).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"notification_id", "user_id", "type", "triggered_by_user_id",
						"pin_id", "comment_id", "message_id", "created_at"}))
			},
			expectedResult: entity.Notification{},
			expectedErr:    sql.ErrNoRows,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			result, err := repo.GetNotificationById(CtxWithRequestId, tc.notificationId)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestDBRepository_ReadNotification(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name           string
		notificationId entity.NotificationID
		setupMocks     func()
		expectedErr    error
	}{
		{
			name:           "OK test",
			notificationId: 1,
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateNotificationStatus).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name:           "Error test",
			notificationId: 1,
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateNotificationStatus).
					WithArgs(1).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name:           "No rows affected",
			notificationId: 1,
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateNotificationStatus).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.ReadNotification(CtxWithRequestId, tc.notificationId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestDBRepository_ReadAllNotifications(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name        string
		userId      entity.UserID
		setupMocks  func()
		expectedErr error
	}{
		{
			name:   "OK test",
			userId: 1,
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateAllNotificationsStatuses).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name:   "Error test",
			userId: 1,
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateAllNotificationsStatuses).
					WithArgs(1).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name:   "No rows affected",
			userId: 1,
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateAllNotificationsStatuses).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.ReadAllNotifications(CtxWithRequestId, tc.userId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
