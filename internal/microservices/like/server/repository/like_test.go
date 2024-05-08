package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/microservices/image/proto"
	mock_proto "harmonica/mocks/microservices/image/proto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_SetLike(t *testing.T) {
	tests := []struct {
		name        string
		pinID       entity.PinID
		userID      entity.UserID
		expectedErr error
	}{
		{
			name:        "OK test",
			pinID:       1,
			userID:      1,
			expectedErr: nil,
		},
		{
			name:        "Error test",
			pinID:       2,
			userID:      2,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, ctrl, _, repo := SetupDBMock(t)
			defer ctrl.Finish()
			defer db.Close()
			mock.ExpectExec(QuerySetLike).
				WithArgs(tc.pinID, tc.userID).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.expectedErr)
			err := repo.SetLike(context.Background(), tc.pinID, tc.userID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_ClearLike(t *testing.T) {
	tests := []struct {
		name        string
		pinID       entity.PinID
		userID      entity.UserID
		expectedErr error
	}{
		{
			name:        "OK test",
			pinID:       1,
			userID:      1,
			expectedErr: nil,
		},
		{
			name:        "Error test",
			pinID:       2,
			userID:      2,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, ctrl, _, repo := SetupDBMock(t)
			defer ctrl.Finish()
			defer db.Close()
			mock.ExpectExec(QueryClearLike).
				WithArgs(tc.pinID, tc.userID).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.expectedErr)
			err := repo.ClearLike(context.Background(), tc.pinID, tc.userID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetUsersLiked(t *testing.T) {
	tests := []struct {
		name          string
		pinID         entity.PinID
		limit         int
		expectedUsers entity.UserList
		expectedErr   error
	}{
		{
			name:          "OK test",
			pinID:         1,
			limit:         10,
			expectedUsers: entity.UserList{Users: []entity.UserResponse{{UserId: 1, Nickname: "User1"}, {UserId: 2, Nickname: "User2"}}},
			expectedErr:   nil,
		},
		{
			name:          "Error test",
			pinID:         2,
			limit:         5,
			expectedUsers: entity.UserList{},
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, ctrl, _, repo := SetupDBMock(t)
			defer ctrl.Finish()
			defer db.Close()
			mock.ExpectQuery(QueryGetUsersLiked).
				WithArgs(tc.pinID, tc.limit).
				WillReturnRows(sqlmock.NewRows([]string{"user_id", "nickname"}).
					AddRow(1, "User1").
					AddRow(2, "User2")).
				WillReturnError(tc.expectedErr)
			users, err := repo.GetUsersLiked(context.Background(), tc.pinID, tc.limit)
			assert.Equal(t, tc.expectedUsers, users)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_CheckIsLiked(t *testing.T) {
	tests := []struct {
		name         string
		pinID        entity.PinID
		userID       entity.UserID
		expectedLike bool
		expectedErr  error
	}{
		{
			name:         "OK test 1",
			pinID:        1,
			userID:       1,
			expectedLike: true,
			expectedErr:  nil,
		},
		{
			name:         "OK test 2",
			pinID:        2,
			userID:       2,
			expectedLike: false,
			expectedErr:  nil,
		},
		{
			name:         "Error test",
			pinID:        3,
			userID:       3,
			expectedLike: false,
			expectedErr:  errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, ctrl, _, repo := SetupDBMock(t)
			defer ctrl.Finish()
			defer db.Close()
			mock.ExpectQuery(QueryIsLiked).
				WithArgs(tc.pinID, tc.userID).
				WillReturnRows(sqlmock.NewRows([]string{"is_liked"}).AddRow(tc.expectedLike)).
				WillReturnError(tc.expectedErr)
			liked, err := repo.CheckIsLiked(context.Background(), tc.pinID, tc.userID)
			assert.Equal(t, tc.expectedLike, liked)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetFavorites(t *testing.T) {
	tests := []struct {
		name         string
		userID       entity.UserID
		expectedPins entity.FeedPins
		error1       error
		error2       error
		expectedErr  error
	}{
		{
			name:         "OK test",
			userID:       1,
			expectedPins: entity.FeedPins{Pins: []entity.FeedPinResponse{{PinId: 1}}},
		},
		{
			name:         "Error test 1",
			userID:       2,
			expectedPins: entity.FeedPins{},
			error1:       errors.New("database error"),
			expectedErr:  errors.New("database error"),
		},
		{
			name:         "Error test 2",
			userID:       2,
			expectedPins: entity.FeedPins{},
			error2:       errors.New("database error"),
			expectedErr:  errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, ctrl, imageClient, repo := SetupDBMock(t)
			defer ctrl.Finish()
			defer db.Close()
			mock.ExpectQuery(QueryGetFavorites).WithArgs(tc.userID, 10, 10).
				WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1)).
				WillReturnError(tc.error1)
			imageClient.EXPECT().GetImageBounds(context.Background(), &proto.GetImageBoundsRequest{}).
				Return(&proto.GetImageBoundsResponse{}, tc.error2).MaxTimes(1)
			pins, err := repo.GetFavorites(context.Background(), tc.userID, 10, 10)
			assert.Equal(t, tc.expectedPins, pins)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func SetupDBMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gomock.Controller, *mock_proto.MockImageClient, *DBRepository) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	ctrl := gomock.NewController(t)
	imageClient := mock_proto.NewMockImageClient(ctrl)
	repo := NewDBRepository(sqlx.NewDb(db, "postgres"), zap.L(), imageClient)
	return db, mock, ctrl, imageClient, repo
}
