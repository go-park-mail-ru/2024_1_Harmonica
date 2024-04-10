package test_repository

import (
	"context"
	"database/sql/driver"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var CtxWithRequestId = context.WithValue(context.Background(), "request_id", "request-id-string")

func TestCreateBoard(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	userId := entity.UserID(1)
	board := entity.Board{}

	// Good test
	mock.ExpectBegin()
	mock.ExpectQuery(repository.QueryCreateBoard).WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
	mock.ExpectExec(repository.QueryInsertBoardAuthor).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
	res, err := r.CreateBoard(CtxWithRequestId, board, userId)
	assert.Equal(t, entity.Board{BoardID: entity.BoardID(1)}, res)
	assert.Equal(t, nil, err)

	// Bad test 1
	mock.ExpectBegin()
	mock.ExpectQuery(repository.QueryCreateBoard).WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
	mock.ExpectExec(repository.QueryInsertBoardAuthor).WillReturnError(errs.ErrDBInternal)
	mock.ExpectRollback()
	res, err = r.CreateBoard(CtxWithRequestId, board, userId)
	assert.Equal(t, entity.Board{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)

	// Bad test 2
	mock.ExpectBegin()
	mock.ExpectQuery(repository.QueryCreateBoard).WillReturnError(errs.ErrDBInternal)
	res, err = r.CreateBoard(CtxWithRequestId, board, userId)
	assert.Equal(t, entity.Board{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestGetBoardById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	boardId := entity.BoardID(1)

	// Good test
	mock.ExpectQuery(repository.QueryGetBoardById).WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
	res, err := r.GetBoardById(CtxWithRequestId, boardId)
	assert.Equal(t, entity.Board{BoardID: entity.BoardID(1)}, res)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectQuery(repository.QueryGetBoardById).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetBoardById(CtxWithRequestId, boardId)
	assert.Equal(t, entity.Board{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestGetBoardAuthors(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	boardId := entity.BoardID(1)

	// Good test
	mock.ExpectQuery(repository.QueryGetBoardAuthors).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
	res, err := r.GetBoardAuthors(CtxWithRequestId, boardId)
	assert.Equal(t, []entity.BoardAuthor{{UserId: entity.UserID(1)}}, res)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectQuery(repository.QueryGetBoardAuthors).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetBoardAuthors(CtxWithRequestId, boardId)
	assert.Equal(t, []entity.BoardAuthor{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestGetBoardPins(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	boardId := entity.BoardID(1)
	limit := 10
	offset := 10

	// Good test
	mock.ExpectQuery(repository.QueryGetBoardPins).WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
	res, err := r.GetBoardPins(CtxWithRequestId, boardId, limit, offset)
	assert.Equal(t, []entity.BoardPinResponse{{PinId: entity.PinID(1)}}, res)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectQuery(repository.QueryGetBoardPins).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetBoardPins(CtxWithRequestId, boardId, limit, offset)
	assert.Equal(t, []entity.BoardPinResponse{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestUpdateBoard(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	board := entity.Board{}

	// Good test
	mock.ExpectQuery(repository.QueryUpdateBoard).WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
	res, err := r.UpdateBoard(CtxWithRequestId, board)
	assert.Equal(t, entity.Board{BoardID: entity.BoardID(1)}, res)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectQuery(repository.QueryUpdateBoard).WillReturnError(errs.ErrDBInternal)
	res, err = r.UpdateBoard(CtxWithRequestId, board)
	assert.Equal(t, entity.Board{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}
func TestAddPinToBoard(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	boardId := entity.BoardID(1)
	pinId := entity.PinID(1)

	// Good test
	mock.ExpectExec(repository.QueryAddPinToBoard).
		WithArgs(boardId, pinId).
		WillReturnResult(driver.ResultNoRows)
	err = r.AddPinToBoard(CtxWithRequestId, boardId, pinId)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectExec(repository.QueryAddPinToBoard).
		WithArgs(boardId, pinId).
		WillReturnError(errs.ErrDBInternal)
	err = r.AddPinToBoard(CtxWithRequestId, boardId, pinId)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestDeletePinFromBoard(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	boardId := entity.BoardID(1)
	pinId := entity.PinID(1)

	// Good test
	mock.ExpectExec(repository.QueryDeletePinFromBoard).
		WithArgs(boardId, pinId).
		WillReturnResult(driver.ResultNoRows)
	err = r.DeletePinFromBoard(CtxWithRequestId, boardId, pinId)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectExec(repository.QueryDeletePinFromBoard).
		WithArgs(boardId, pinId).
		WillReturnError(errs.ErrDBInternal)
	err = r.DeletePinFromBoard(CtxWithRequestId, boardId, pinId)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestDeleteBoard(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	boardId := entity.BoardID(1)

	// Good test
	mock.ExpectExec(repository.QueryDeleteBoard).
		WithArgs(boardId).
		WillReturnResult(driver.ResultNoRows)
	err = r.DeleteBoard(CtxWithRequestId, boardId)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectExec(repository.QueryDeleteBoard).
		WithArgs(boardId).
		WillReturnError(errs.ErrDBInternal)
	err = r.DeleteBoard(CtxWithRequestId, boardId)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestGetUserBoards(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	authorId := entity.UserID(1)
	limit := 10
	offset := 10

	// Good test
	mock.ExpectQuery(repository.QueryGetUserBoards).
		WithArgs(authorId, limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"board_id"}).AddRow(1))
	res, err := r.GetUserBoards(CtxWithRequestId, authorId, limit, offset)
	assert.Equal(t, entity.UserBoards{Boards: []entity.Board{{BoardID: entity.BoardID(1)}}}, res)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectQuery(repository.QueryGetUserBoards).
		WithArgs(authorId, limit, offset).
		WillReturnError(errs.ErrDBInternal)
	res, err = r.GetUserBoards(CtxWithRequestId, authorId, limit, offset)
	assert.Equal(t, entity.UserBoards{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestCheckBoardAuthorExistence(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	userId := entity.UserID(1)
	boardId := entity.BoardID(1)

	// Good test
	mock.ExpectQuery(repository.QueryCheckBoardAuthorExistence).
		WithArgs(userId, boardId).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(true))
	res, err := r.CheckBoardAuthorExistence(CtxWithRequestId, userId, boardId)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	// Error test
	mock.ExpectQuery(repository.QueryCheckBoardAuthorExistence).
		WithArgs(userId, boardId).
		WillReturnError(errs.ErrDBInternal)
	res, err = r.CheckBoardAuthorExistence(CtxWithRequestId, userId, boardId)
	assert.Equal(t, false, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}
