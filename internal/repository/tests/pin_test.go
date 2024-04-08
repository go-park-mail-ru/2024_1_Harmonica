package test_repository

import (
	"context"
	"database/sql/driver"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetFeedPins(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	limit := 10
	offset := 10
	// Good test
	mock.ExpectQuery(repository.QueryGetPinsFeed).WillReturnRows(sqlmock.NewRows([]string{}))
	res, err := r.GetFeedPins(context.Background(), limit, offset)
	assert.Equal(t, entity.FeedPins{}, res)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectQuery(repository.QueryGetPinsFeed).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetFeedPins(context.Background(), limit, offset)
	assert.Equal(t, entity.FeedPins{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestGetUserPins(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	limit := 10
	offset := 10
	authorId := entity.UserID(1)

	// Good test
	mock.ExpectQuery(repository.QueryGetUserPins).WillReturnRows(sqlmock.NewRows([]string{}))
	res, err := r.GetUserPins(context.Background(), authorId, limit, offset)
	assert.Equal(t, entity.UserPins{}, res)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectQuery(repository.QueryGetUserPins).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetUserPins(context.Background(), authorId, limit, offset)
	assert.Equal(t, entity.UserPins{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestGetPinById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)

	pinId := entity.PinID(1)
	// Good test
	mock.ExpectQuery(repository.QueryGetPinById).
		WillReturnRows(sqlmock.NewRows([]string{"content_url"}).
			AddRow("https://img/image.png.svg"))
	res, err := r.GetPinById(context.Background(), pinId)
	assert.Equal(t, entity.PinPageResponse{ContentUrl: "https://img/image.png.svg"}, res)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectQuery(repository.QueryGetPinById).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetPinById(context.Background(), pinId)
	assert.Equal(t, entity.PinPageResponse{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestCreatePin(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)

	pin := entity.Pin{Title: "Title"}
	excpectedID := entity.PinID(13)
	// Good test
	mock.ExpectQuery(repository.QueryCreatePin).
		WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow("13"))
	res, err := r.CreatePin(context.Background(), pin)
	assert.Equal(t, excpectedID, res)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectQuery(repository.QueryCreatePin).WillReturnError(errs.ErrDBInternal)
	res, err = r.CreatePin(context.Background(), pin)
	assert.Equal(t, entity.PinID(0), res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestUpdatePin(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)

	pin := entity.Pin{Title: "Title"}
	// Good test
	mock.ExpectExec(repository.QueryUpdatePin).WillReturnResult(driver.ResultNoRows)
	err = r.UpdatePin(context.Background(), pin)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectExec(repository.QueryUpdatePin).WillReturnError(errs.ErrDBInternal)
	err = r.UpdatePin(context.Background(), pin)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestDeletePin(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	pinId := entity.PinID(1)

	// Good test
	mock.ExpectExec(repository.QueryDeletePin).WillReturnResult(driver.ResultNoRows)
	err = r.DeletePin(context.Background(), pinId)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectExec(repository.QueryDeletePin).WillReturnError(errs.ErrDBInternal)
	err = r.DeletePin(context.Background(), pinId)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestCheckPinExistence(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	pinId := entity.PinID(1)

	// Good test
	mock.ExpectQuery(repository.QueryCheckPinExistence).WithArgs(entity.PinID(1)).WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("true"))
	res, err := r.CheckPinExistence(context.Background(), pinId)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectQuery(repository.QueryCheckPinExistence).WithArgs(entity.PinID(1)).WillReturnError(errs.ErrDBInternal)
	res, err = r.CheckPinExistence(context.Background(), pinId)
	assert.Equal(t, false, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}
