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

func TestSetLike(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	pinId := entity.PinID(1)
	userId := entity.UserID(1)

	// Good test
	mock.ExpectExec(repository.QuerySetLike).WillReturnResult(driver.ResultNoRows)
	err = r.SetLike(context.Background(), pinId, userId)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectExec(repository.QuerySetLike).WillReturnError(errs.ErrDBInternal)
	err = r.SetLike(context.Background(), pinId, userId)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestClearLike(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	pinId := entity.PinID(1)
	userId := entity.UserID(1)

	// Good test
	mock.ExpectExec(repository.QueryClearLike).WillReturnResult(driver.ResultNoRows)
	err = r.ClearLike(context.Background(), pinId, userId)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectExec(repository.QueryClearLike).WillReturnError(errs.ErrDBInternal)
	err = r.ClearLike(context.Background(), pinId, userId)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestGetUsersLiked(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	pinId := entity.PinID(1)
	limit := 20

	// Good test
	mock.ExpectQuery(repository.QueryGetUsersLiked).WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("email"))
	res, err := r.GetUsersLiked(context.Background(), pinId, limit)
	expectedRes := entity.UserList{Users: []entity.UserResponse{{Email: "email"}}}
	assert.Equal(t, expectedRes, res)
	assert.Equal(t, nil, err)
	// Error test 1
	mock.ExpectQuery(repository.QueryGetUsersLiked).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetUsersLiked(context.Background(), pinId, limit)
	assert.Equal(t, entity.UserList{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestCheckIsLiked(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil)
	pinId := entity.PinID(1)
	userId := entity.UserID(1)

	// Good test
	mock.ExpectQuery(repository.QueryIsLiked).WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("true"))
	res, err := r.CheckIsLiked(context.Background(), pinId, userId)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)
	// Error test 1
	mock.ExpectQuery(repository.QueryIsLiked).WillReturnError(errs.ErrDBInternal)
	res, err = r.CheckIsLiked(context.Background(), pinId, userId)
	assert.Equal(t, false, res)
	assert.Equal(t, errs.ErrDBInternal, err)
}
