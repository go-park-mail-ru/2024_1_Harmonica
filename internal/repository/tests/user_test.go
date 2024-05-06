package test_repository

/*
import (
	"database/sql/driver"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	email := "emaiiil@godot.com"
	// Good test
	mock.ExpectQuery(repository.QueryGetUserByEmail).WillReturnRows(sqlmock.NewRows([]string{"nickname"}).AddRow("Nikoin"))
	res, err := r.GetUserByEmail(CtxWithRequestId, email)
	assert.Equal(t, entity.User{Nickname: "Nikoin"}, res)
	assert.Equal(t, nil, err)
	// Error test 1
	mock.ExpectQuery(repository.QueryGetUserByEmail).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetUserByEmail(CtxWithRequestId, email)
	assert.Equal(t, entity.User{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
	// Error test 2
	mock.ExpectQuery(repository.QueryGetUserByEmail).WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
	res, err = r.GetUserByEmail(CtxWithRequestId, email)
	assert.Equal(t, entity.User{}, res)
	assert.Error(t, err)
}

func TestGetUserByNickname(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	CtxWithRequestId := context.Background()
	CtxWithRequestId = context.WithValue(CtxWithRequestId, "request_id", "1935")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	nickname := "Nikoin"
	// Good test
	mock.ExpectQuery(repository.QueryGetUserByNickname).WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("emaiiil@godot.com"))
	res, err := r.GetUserByNickname(CtxWithRequestId, nickname)
	assert.Equal(t, entity.User{Email: "emaiiil@godot.com"}, res)
	assert.Equal(t, nil, err)
	// Error test 1
	mock.ExpectQuery(repository.QueryGetUserByNickname).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetUserByNickname(CtxWithRequestId, nickname)
	assert.Equal(t, entity.User{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
	// Error test 2
	mock.ExpectQuery(repository.QueryGetUserByNickname).WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
	res, err = r.GetUserByNickname(CtxWithRequestId, nickname)
	assert.Equal(t, entity.User{}, res)
	assert.Error(t, err)
}

func TestGetUserById(t *testing.T) {
	CtxWithRequestId := context.Background()
	CtxWithRequestId = context.WithValue(CtxWithRequestId, "request_id", "1935")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	id := entity.UserID(1)
	// Good test
	mock.ExpectQuery(repository.QueryGetUserById).WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("emaiiil@godot.com"))
	res, err := r.GetUserById(CtxWithRequestId, id)
	assert.Equal(t, entity.User{Email: "emaiiil@godot.com"}, res)
	assert.Equal(t, nil, err)
	// Error test 1
	mock.ExpectQuery(repository.QueryGetUserById).WillReturnError(errs.ErrDBInternal)
	res, err = r.GetUserById(CtxWithRequestId, id)
	assert.Equal(t, entity.User{}, res)
	assert.Equal(t, errs.ErrDBInternal, err)
	// Error test 2
	mock.ExpectQuery(repository.QueryGetUserById).WillReturnRows(sqlmock.NewRows([]string{"err"}).AddRow("err"))
	res, err = r.GetUserById(CtxWithRequestId, id)
	assert.Equal(t, entity.User{}, res)
	assert.Error(t, err)
}

func TestRegisterUser(t *testing.T) {
	CtxWithRequestId := context.Background()
	CtxWithRequestId = context.WithValue(CtxWithRequestId, "request_id", "1935")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	user := entity.User{
		Nickname:  "nickname",
		Password:  "hashstring",
		AvatarURL: "htttp10://img.png",
		Email:     "email@com",
	}

	// Good test
	mock.ExpectExec(repository.QueryRegisterUser).WithArgs(user.Email, user.Nickname, user.Password).
		WillReturnResult(driver.ResultNoRows)
	err = r.RegisterUser(CtxWithRequestId, user)
	assert.Equal(t, nil, err)
	// Error test
	mock.ExpectExec(repository.QueryRegisterUser).WithArgs(user.Email, user.Nickname, user.Password).
		WillReturnError(errs.ErrDBInternal)
	err = r.RegisterUser(CtxWithRequestId, user)
	assert.Equal(t, errs.ErrDBInternal, err)
}

func TestUpdateUser(t *testing.T) {
	CtxWithRequestId := context.Background()
	CtxWithRequestId = context.WithValue(CtxWithRequestId, "request_id", "1935")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r := repository.NewDBRepository(sqlx.NewDb(db, "postgres"), nil, zap.L())
	user := entity.User{
		Nickname:  "nickname",
		Password:  "hashstring",
		AvatarURL: "htttp10://img.png",
		Email:     "email@com",
		UserID:    entity.UserID(1),
	}

	// Good test
	mock.ExpectExec(repository.QueryUpdateUserNickname).WithArgs(user.UserID, user.Nickname).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(repository.QueryUpdateUserPassword).WithArgs(user.UserID, user.Password).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(repository.QueryUpdateUserAvatar).WithArgs(user.UserID, user.AvatarURL).
		WillReturnResult(driver.ResultNoRows)

	err = r.UpdateUser(CtxWithRequestId, user)
	assert.Equal(t, nil, err)
	// Error test 1
	mock.ExpectExec(repository.QueryUpdateUserNickname).WithArgs(user.UserID, user.Nickname).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(repository.QueryUpdateUserPassword).WithArgs(user.UserID, user.Password).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(repository.QueryUpdateUserAvatar).WithArgs(user.UserID, user.AvatarURL).
		WillReturnError(errs.ErrDBInternal)

	err = r.UpdateUser(CtxWithRequestId, user)
	assert.Equal(t, errs.ErrDBInternal, err)

	// Error test 2
	mock.ExpectExec(repository.QueryUpdateUserNickname).WithArgs(user.UserID, user.Nickname).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(repository.QueryUpdateUserPassword).WithArgs(user.UserID, user.Password).
		WillReturnError(errs.ErrDBInternal)

	err = r.UpdateUser(CtxWithRequestId, user)
	assert.Equal(t, errs.ErrDBInternal, err)

	// Error test 3
	mock.ExpectExec(repository.QueryUpdateUserNickname).WithArgs(user.UserID, user.Nickname).
		WillReturnError(errs.ErrDBInternal)

	err = r.UpdateUser(CtxWithRequestId, user)
	assert.Equal(t, errs.ErrDBInternal, err)
}
*/
