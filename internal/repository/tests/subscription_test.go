package test_repository

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/repository"
	"testing"
)

func TestRepository_AddSubscriptionToUser(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name        string
		userID      entity.UserID
		subscribeID entity.UserID
		setupMocks  func()
		expectedErr error
	}{
		{
			name:        "OK test",
			userID:      entity.UserID(1),
			subscribeID: entity.UserID(2),
			setupMocks: func() {
				mock.ExpectExec(repository.QueryAddSubscriptionToUser).
					WithArgs(entity.UserID(1), entity.UserID(2)).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name:        "Error test",
			userID:      entity.UserID(1),
			subscribeID: entity.UserID(2),
			setupMocks: func() {
				mock.ExpectExec(repository.QueryAddSubscriptionToUser).
					WithArgs(entity.UserID(1), entity.UserID(2)).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.AddSubscriptionToUser(CtxWithRequestId, tc.userID, tc.subscribeID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_DeleteSubscriptionToUser(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name          string
		userID        entity.UserID
		unsubscribeID entity.UserID
		setupMocks    func()
		expectedErr   error
	}{
		{
			name:          "OK test",
			userID:        entity.UserID(1),
			unsubscribeID: entity.UserID(2),
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeleteSubscriptionToUser).
					WithArgs(entity.UserID(1), entity.UserID(2)).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name:          "Error test",
			userID:        entity.UserID(1),
			unsubscribeID: entity.UserID(2),
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeleteSubscriptionToUser).
					WithArgs(entity.UserID(1), entity.UserID(2)).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.DeleteSubscriptionToUser(CtxWithRequestId, tc.userID, tc.unsubscribeID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetSubscriptionsInfo(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name            string
		userToGetInfoID entity.UserID
		userID          entity.UserID
		setupMocks      func()
		expectedInfo    entity.UserProfileResponse
		expectedErr     error
	}{
		{
			name:            "OK test",
			userToGetInfoID: entity.UserID(1),
			userID:          entity.UserID(2),
			setupMocks: func() {
				mock.ExpectQuery(repository.GetSubscriptionsInfo).
					WithArgs(entity.UserID(1), entity.UserID(2)).
					WillReturnRows(sqlmock.NewRows([]string{"subscribers_count", "subscriptions_count", "is_subscribed"}).
						AddRow(10, 5, true))
			},
			expectedInfo: entity.UserProfileResponse{
				SubscribersCount:   10,
				SubscriptionsCount: 5,
				IsSubscribed:       true,
			},
			expectedErr: nil,
		},
		{
			name:            "Error test",
			userToGetInfoID: entity.UserID(1),
			userID:          entity.UserID(2),
			setupMocks: func() {
				mock.ExpectQuery(repository.GetSubscriptionsInfo).
					WithArgs(entity.UserID(1), entity.UserID(2)).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedInfo: entity.UserProfileResponse{},
			expectedErr:  errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			info, err := repo.GetSubscriptionsInfo(CtxWithRequestId, tc.userToGetInfoID, tc.userID)
			assert.Equal(t, tc.expectedInfo, info)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetUserSubscribers(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name                string
		userID              entity.UserID
		setupMocks          func()
		expectedSubscribers entity.UserSubscribers
		expectedErr         error
	}{
		{
			name:   "OK test",
			userID: entity.UserID(1),
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserSubscribers).
					WithArgs(entity.UserID(1)).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "nickname", "avatar_url"}).
						AddRow(1, "user", "url"))
			},
			expectedSubscribers: entity.UserSubscribers{
				Subscribers: []entity.UserResponse{
					{UserId: 1, Nickname: "user", AvatarURL: "url"},
				},
			},
			expectedErr: nil,
		},
		{
			name:   "Error test",
			userID: entity.UserID(1),
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserSubscribers).
					WithArgs(entity.UserID(1)).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedSubscribers: entity.UserSubscribers{},
			expectedErr:         errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			subscribers, err := repo.GetUserSubscribers(CtxWithRequestId, tc.userID)
			assert.Equal(t, tc.expectedSubscribers, subscribers)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetUserSubscriptions(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name                  string
		userID                entity.UserID
		setupMocks            func()
		expectedSubscriptions entity.UserSubscriptions
		expectedErr           error
	}{
		{
			name:   "OK test",
			userID: entity.UserID(1),
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserSubscriptions).
					WithArgs(entity.UserID(1)).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "nickname", "avatar_url"}).
						AddRow(1, "user", "url"))
			},
			expectedSubscriptions: entity.UserSubscriptions{
				Subscriptions: []entity.UserResponse{
					{UserId: 1, Nickname: "user", AvatarURL: "url"},
				},
			},
			expectedErr: nil,
		},
		{
			name:   "Error test",
			userID: entity.UserID(1),
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserSubscriptions).
					WithArgs(entity.UserID(1)).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedSubscriptions: entity.UserSubscriptions{},
			expectedErr:           errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			subscriptions, err := repo.GetUserSubscriptions(CtxWithRequestId, tc.userID)
			assert.Equal(t, tc.expectedSubscriptions, subscriptions)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
