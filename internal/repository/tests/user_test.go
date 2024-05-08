package test_repository

import (
	"database/sql/driver"
	"errors"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/microservices/image/proto"
	"harmonica/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
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
				mock.ExpectQuery(repository.QueryGetUserByEmail).
					WillReturnRows(sqlmock.NewRows([]string{"nickname", "avatar_url"}).
						AddRow("nickname", "url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{Dx: 1}, nil).MaxTimes(1)
			},
			expectedUser: entity.User{Nickname: "nickname", AvatarURL: "url", AvatarDX: 1},
			expectedErr:  nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserByEmail).WillReturnError(errs.ErrDBInternal)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserByEmail).
					WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
			},
			expectedUser: entity.User{},
			expectedErr:  errors.New("missing destination name err in *entity.User"),
		},
		{
			name: "Error test 3",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserByEmail).
					WillReturnRows(sqlmock.NewRows([]string{"nickname", "avatar_url"}).
						AddRow("nickname", "url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{Dx: 1}, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
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

func TestGetUserByNickname(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
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
				mock.ExpectQuery(repository.QueryGetUserByNickname).
					WillReturnRows(sqlmock.NewRows([]string{"nickname", "avatar_url"}).
						AddRow("nickname", "url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{Dx: 1, Dy: 2}, nil).MaxTimes(1)
			},
			expectedUser: entity.User{Nickname: "nickname", AvatarURL: "url", AvatarDX: 1, AvatarDY: 2},
			expectedErr:  nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserByNickname).WillReturnError(errs.ErrDBInternal)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserByNickname).
					WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
			},
			expectedUser: entity.User{},
			expectedErr:  errors.New("missing destination name err in *entity.User"),
		},
		{
			name: "Error test 3",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserByNickname).
					WillReturnRows(sqlmock.NewRows([]string{"nickname", "avatar_url"}).
						AddRow("nickname", "url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetUserByNickname(CtxWithRequestId, "nickname")
			assert.Equal(t, tc.expectedUser, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetUserById(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
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
				mock.ExpectQuery(repository.QueryGetUserById).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "avatar_url"}).
						AddRow(1, "url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{Dx: 1, Dy: 2}, nil).MaxTimes(1)
			},
			expectedUser: entity.User{UserID: 1, AvatarURL: "url", AvatarDX: 1, AvatarDY: 2},
			expectedErr:  nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserById).WillReturnError(errs.ErrDBInternal)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserById).
					WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
			},
			expectedUser: entity.User{},
			expectedErr:  errors.New("missing destination name err in *entity.User"),
		},
		{
			name: "Error test 3",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserById).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "avatar_url"}).
						AddRow(1, "url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedUser: entity.User{},
			expectedErr:  errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetUserById(CtxWithRequestId, 1)
			assert.Equal(t, tc.expectedUser, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRegisterUser(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Good test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryRegisterUser).
					WithArgs("email@com", "nickname", "hashPassword").
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryRegisterUser).
					WithArgs("email@com", "nickname", "hashPassword").
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.RegisterUser(CtxWithRequestId, entity.User{
				Nickname:  "nickname",
				Password:  "hashPassword",
				AvatarURL: "h://img.png",
				Email:     "email@com",
			})
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Good test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateUserNickname).
					WithArgs(entity.UserID(1), "nickname").
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(repository.QueryUpdateUserPassword).
					WithArgs(entity.UserID(1), "hashPassword").
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(repository.QueryUpdateUserAvatar).
					WithArgs(entity.UserID(1), "h://img.png").
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateUserNickname).
					WithArgs(entity.UserID(1), "nickname").
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(repository.QueryUpdateUserPassword).
					WithArgs(entity.UserID(1), "hashPassword").
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(repository.QueryUpdateUserAvatar).
					WithArgs(entity.UserID(1), "h://img.png").
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateUserNickname).
					WithArgs(entity.UserID(1), "nickname").
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec(repository.QueryUpdateUserPassword).
					WithArgs(entity.UserID(1), "hashPassword").
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test 3",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdateUserNickname).
					WithArgs(entity.UserID(1), "nickname").
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.UpdateUser(CtxWithRequestId, entity.User{
				Nickname:  "nickname",
				Password:  "hashPassword",
				AvatarURL: "h://img.png",
				Email:     "email@com",
				UserID:    entity.UserID(1),
			})
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
