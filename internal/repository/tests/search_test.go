package test_repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/repository"
	"testing"
)

const SearchLimit = 10

func TestRepository_SearchForUsers(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name          string
		query         string
		setupMocks    func()
		expectedUsers []entity.SearchUser
		expectedErr   error
	}{
		{
			name:  "OK test 1",
			query: "searchQuery",
			setupMocks: func() {
				mock.ExpectQuery(repository.QuerySearchForUser).
					WithArgs("searchQuery", SearchLimit).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "nickname", "avatar_url", "subs"}).
						AddRow(1, "user1", "avatar1", 10))
			},
			expectedUsers: []entity.SearchUser{
				{UserId: entity.UserID(1), Nickname: "user1", AvatarURL: "avatar1", SubsCount: 10},
			},
			expectedErr: nil,
		},
		{
			name:  "Error test 1",
			query: "searchQuery",
			setupMocks: func() {
				mock.ExpectQuery(repository.QuerySearchForUser).
					WithArgs("searchQuery", SearchLimit).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedUsers: []entity.SearchUser{},
			expectedErr:   errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			users, err := repo.SearchForUsers(CtxWithRequestId, tc.query)
			assert.Equal(t, tc.expectedUsers, users)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_SearchForPins(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name         string
		query        string
		setupMocks   func()
		expectedPins []entity.SearchPin
		expectedErr  error
	}{
		{
			name:  "OK test 1",
			query: "searchQuery",
			setupMocks: func() {
				mock.ExpectQuery(repository.QuerySearchForPin).
					WithArgs("searchQuery", SearchLimit).
					WillReturnRows(sqlmock.NewRows([]string{"pin_id", "title", "content_url"}).
						AddRow(1, "pin1", "url1").
						AddRow(2, "pin2", "url2"))
			},
			expectedPins: []entity.SearchPin{
				{PinId: entity.PinID(1), Title: "pin1", ContentURL: "url1"},
				{PinId: entity.PinID(2), Title: "pin2", ContentURL: "url2"},
			},
			expectedErr: nil,
		},
		{
			name:  "Error test 1",
			query: "searchQuery",
			setupMocks: func() {
				mock.ExpectQuery(repository.QuerySearchForPin).
					WithArgs("searchQuery", SearchLimit).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedPins: []entity.SearchPin{},
			expectedErr:  errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			pins, err := repo.SearchForPins(CtxWithRequestId, tc.query)
			assert.Equal(t, tc.expectedPins, pins)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_SearchForBoards(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name           string
		query          string
		setupMocks     func()
		expectedBoards []entity.SearchBoard
		expectedErr    error
	}{
		{
			name:  "OK test 1",
			query: "searchQuery",
			setupMocks: func() {
				mock.ExpectQuery(repository.QuerySearchForBoard).
					WithArgs("searchQuery", SearchLimit).
					WillReturnRows(sqlmock.NewRows([]string{"board_id", "title", "cover_url"}).
						AddRow(1, "board1", "cover1").
						AddRow(2, "board2", "cover2"))
			},
			expectedBoards: []entity.SearchBoard{
				{BoardId: entity.BoardID(1), Title: "board1", CoverURL: "cover1"},
				{BoardId: entity.BoardID(2), Title: "board2", CoverURL: "cover2"},
			},
			expectedErr: nil,
		},
		{
			name:  "Error test 1",
			query: "searchQuery",
			setupMocks: func() {
				mock.ExpectQuery(repository.QuerySearchForBoard).
					WithArgs("searchQuery", SearchLimit).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedBoards: []entity.SearchBoard{},
			expectedErr:    errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			boards, err := repo.SearchForBoards(CtxWithRequestId, tc.query)
			assert.Equal(t, tc.expectedBoards, boards)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
