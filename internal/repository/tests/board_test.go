package test_repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/microservices/image/proto"
	"harmonica/internal/repository"
	mock_proto "harmonica/mocks/microservices/image/proto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var CtxWithRequestId = context.WithValue(context.Background(), "request_id", "request-id-string")

func TestRepository_CreateBoard(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()
	tests := []struct {
		name        string
		setupMocks  func()
		expectedRes entity.Board
		expectedErr error
	}{
		{
			name: "OK test case 1",
			setupMocks: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(repository.QueryCreateBoard).
					WillReturnRows(sqlmock.NewRows([]string{"board_id", "cover_url"}).AddRow(1, "url"))
				mock.ExpectExec(repository.QueryInsertBoardAuthor).WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{Dx: 1, Dy: 2}, nil)
			},
			expectedRes: entity.Board{BoardID: 1, CoverURL: "url", CoverDX: 1, CoverDY: 2},
			expectedErr: nil,
		},
		{
			name: "Error test case 1",
			setupMocks: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(repository.QueryCreateBoard).WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil)
				mock.ExpectExec(repository.QueryInsertBoardAuthor).WillReturnError(errs.ErrDBInternal)
				mock.ExpectRollback()
			},
			expectedRes: entity.Board{},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test case 2",
			setupMocks: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(repository.QueryCreateBoard).WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal)
				mock.ExpectRollback()
			},
			expectedRes: entity.Board{},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test case 3",
			setupMocks: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(repository.QueryCreateBoard).WillReturnError(errs.ErrDBInternal)
			},
			expectedRes: entity.Board{},
			expectedErr: errs.ErrDBInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			res, err := repo.CreateBoard(CtxWithRequestId, entity.Board{}, 1)
			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestRepository_GetBoardById(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()
	tests := []struct {
		name        string
		boardID     entity.BoardID
		setupMocks  func()
		expectedRes entity.Board
		expectedErr error
	}{
		{
			name:    "OK test case 1",
			boardID: 1,
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardById).
					WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{Dx: 1, Dy: 2}, nil).MaxTimes(1)
			},
			expectedRes: entity.Board{BoardID: 1, CoverDX: 1, CoverDY: 2},
			expectedErr: nil,
		},
		{
			name:    "Error test case 1",
			boardID: 1,
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardById).
					WillReturnError(errs.ErrDBInternal)
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil).MaxTimes(1)
			},
			expectedRes: entity.Board{},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name:    "Error test case 2",
			boardID: 1,
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardById).
					WillReturnError(errs.ErrDBInternal)
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedRes: entity.Board{},
			expectedErr: errs.ErrDBInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			res, err := repo.GetBoardById(CtxWithRequestId, tt.boardID)
			assert.Equal(t, tt.expectedRes, res)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestRepository_GetBoardAuthors(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()
	tests := []struct {
		name        string
		setupMocks  func()
		expectedRes []entity.BoardAuthor
		expectedErr error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardAuthors).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil).MaxTimes(1)
			},
			expectedRes: []entity.BoardAuthor{{UserId: entity.UserID(1)}},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardAuthors).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedRes: []entity.BoardAuthor{},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardAuthors).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedRes: []entity.BoardAuthor{},
			expectedErr: errs.ErrDBInternal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetBoardAuthors(CtxWithRequestId, 1)
			assert.Equal(t, tc.expectedRes, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetBoardPins(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name        string
		setupMocks  func()
		expectedRes []entity.BoardPinResponse
		expectedErr error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardPins).WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil).MaxTimes(1)
			},
			expectedRes: []entity.BoardPinResponse{{PinId: entity.PinID(1)}},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardPins).WillReturnError(errs.ErrDBInternal)
			},
			expectedRes: []entity.BoardPinResponse{},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetBoardPins).WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedRes: []entity.BoardPinResponse{},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetBoardPins(CtxWithRequestId, 1, 10, 10)
			assert.Equal(t, tc.expectedRes, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_UpdateBoard(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name          string
		setupMocks    func()
		expectedBoard entity.Board
		expectedErr   error
	}{
		{
			name: "OK test",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryUpdateBoard).
					WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil).MaxTimes(1)
			},
			expectedBoard: entity.Board{BoardID: entity.BoardID(1)},
			expectedErr:   nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryUpdateBoard).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedBoard: entity.Board{},
			expectedErr:   errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryUpdateBoard).
					WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedBoard: entity.Board{},
			expectedErr:   errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.UpdateBoard(CtxWithRequestId, entity.Board{})
			assert.Equal(t, tc.expectedBoard, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_AddPinToBoard(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	boardId := entity.BoardID(1)
	pinId := entity.PinID(1)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "OK test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryAddPinToBoard).
					WithArgs(boardId, pinId).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryAddPinToBoard).
					WithArgs(boardId, pinId).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.AddPinToBoard(CtxWithRequestId, boardId, pinId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_DeletePinFromBoard(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	boardId := entity.BoardID(1)
	pinId := entity.PinID(1)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Good test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeletePinFromBoard).
					WithArgs(boardId, pinId).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeletePinFromBoard).
					WithArgs(boardId, pinId).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.DeletePinFromBoard(CtxWithRequestId, boardId, pinId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_DeleteBoard(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	boardId := entity.BoardID(1)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "Good test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeleteBoard).
					WithArgs(boardId).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeleteBoard).
					WithArgs(boardId).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.DeleteBoard(CtxWithRequestId, boardId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetUserBoards(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	authorID := entity.UserID(1)
	userID := entity.UserID(2)
	limit := 10
	offset := 0

	tests := []struct {
		name           string
		setupMocks     func()
		expectedBoards entity.UserBoards
		expectedErr    error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserBoards).
					WithArgs(authorID, userID, limit, offset).
					WillReturnRows(sqlmock.NewRows([]string{"board_id", "recent_pins"}).
						AddRow(1, "{}").
						AddRow(2, "{url1,url2,url3}"))
			},
			expectedBoards: entity.UserBoards{
				Boards: []entity.UserBoard{
					{BoardID: 1, RecentPinContentUrls: nil},
					{BoardID: 2, RecentPinContentUrls: []string{"url1", "url2", "url3"}},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserBoards).
					WithArgs(authorID, userID, limit, offset).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedBoards: entity.UserBoards{},
			expectedErr:    errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetUserBoards(CtxWithRequestId, authorID, userID, limit, offset)
			assert.Equal(t, tc.expectedBoards, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_CheckBoardAuthorExistence(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	userId := entity.UserID(1)
	boardId := entity.BoardID(1)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedRes bool
		expectedErr error
	}{
		{
			name: "Good test",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCheckBoardAuthorExistence).
					WithArgs(userId, boardId).
					WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(true))
			},
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "Error test",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCheckBoardAuthorExistence).
					WithArgs(userId, boardId).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedRes: false,
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.CheckBoardAuthorExistence(CtxWithRequestId, userId, boardId)
			assert.Equal(t, tc.expectedRes, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func SetupDBMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gomock.Controller, *mock_proto.MockImageClient, *repository.DBRepository) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	ctrl := gomock.NewController(t)
	imageClient := mock_proto.NewMockImageClient(ctrl)
	repo := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), imageClient, zap.L())
	return db, mock, ctrl, imageClient, repo
}
