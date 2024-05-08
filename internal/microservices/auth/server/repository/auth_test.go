package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"testing"
)

var CtxWithRequestId = context.WithValue(context.Background(), "request_id", "request-id-string")

func TestRepository_GetUserByEmail(t *testing.T) {
	db, mock, repo := SetupDBMock(t)
	defer db.Close()

	tests := []struct {
		name         string
		setupMocks   func()
		expectedUser entity.User
		expectedErr  error
	}{
		{
			name: "Good test",
			setupMocks: func() {
				mock.ExpectQuery(QueryGetUserByEmail).
					WillReturnRows(sqlmock.NewRows([]string{"nickname", "avatar_url"}).
						AddRow("nickname", "url"))
			},
			expectedUser: entity.User{Nickname: "nickname", AvatarURL: "url"},
			expectedErr:  nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(QueryGetUserByEmail).WillReturnError(errs.ErrDBInternal)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(QueryGetUserByEmail).
					WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
			},
			expectedUser: entity.User{},
			expectedErr:  errors.New("missing destination name err in *entity.User"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetUserByEmail(CtxWithRequestId, "emaiiil@godot.com")
			assert.Equal(t, tc.expectedUser, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetUserById(t *testing.T) {
	db, mock, repo := SetupDBMock(t)
	defer db.Close()

	tests := configureTests(mock)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetUserById(CtxWithRequestId, 1)
			assert.Equal(t, tc.expectedUser, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func SetupDBMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *DBRepository) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo := NewDBRepository(sqlx.NewDb(db, "postgres"), zap.L())
	return db, mock, repo
}

type testStruct struct {
	name         string
	setupMocks   func()
	expectedUser entity.User
	expectedErr  error
}

func configureTests(mock sqlmock.Sqlmock) []testStruct {
	tests := []testStruct{
		{
			name: "Good test",
			setupMocks: func() {
				mock.ExpectQuery(QueryGetUserById).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "avatar_url"}).
						AddRow(1, "url"))
			},
			expectedUser: entity.User{UserID: 1, AvatarURL: "url"},
			expectedErr:  nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(QueryGetUserById).WillReturnError(errs.ErrDBInternal)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(QueryGetUserById).
					WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
			},
			expectedUser: entity.User{},
			expectedErr:  errors.New("missing destination name err in *entity.User"),
		},
	}
	return tests
}
