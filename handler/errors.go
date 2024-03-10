package handler

import (
	"encoding/json"
	"errors"
	"harmonica/models"
	"log"
	"net/http"
)

var (
	ErrAlreadyAuthorized  = errors.New("already authorized")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrReadCookie         = errors.New("error reading cookie")
	ErrReadingRequestBody = errors.New("error reading request body")
	ErrInvalidInputFormat = errors.New("failed to login/register because of invalid input format")
	ErrHashingPassword    = errors.New("error hashing password")
	ErrUserNotExist       = errors.New("user with this email does not exist (can't authorize)")
	ErrWrongPassword      = errors.New("wrong password (can't authorize)")
	ErrDBUniqueEmail      = errors.New("user with this email already exists (can't register)")
	ErrDBUniqueNickname   = errors.New("user with this nickname already exists (can't register)")
	ErrDBInternal         = errors.New("internal db error")
)

var ErrorCodes = map[error]struct {
	HttpCode  int
	LocalCode int
}{
	ErrAlreadyAuthorized:  {HttpCode: 403, LocalCode: 1},
	ErrUnauthorized:       {HttpCode: 401, LocalCode: 2},
	ErrReadCookie:         {HttpCode: 400, LocalCode: 3},
	ErrReadingRequestBody: {HttpCode: 400, LocalCode: 4},
	ErrInvalidInputFormat: {HttpCode: 400, LocalCode: 5},
	ErrHashingPassword:    {HttpCode: 500, LocalCode: 6},
	ErrUserNotExist:       {HttpCode: 401, LocalCode: 7},
	ErrWrongPassword:      {HttpCode: 401, LocalCode: 8},
	ErrDBUniqueEmail:      {HttpCode: 500, LocalCode: 9},
	ErrDBUniqueNickname:   {HttpCode: 500, LocalCode: 10},
	ErrDBInternal:         {HttpCode: 500, LocalCode: 11},
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	log.Println("ERROR ", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorCodes[err].HttpCode)
	response, _ := json.Marshal(models.ErrorResponse{
		Code:    ErrorCodes[err].LocalCode,
		Message: err.Error(),
	})
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
}

func WriteErrorsListResponse(w http.ResponseWriter, errors ...error) {
	var list []models.ErrorResponse
	for _, err := range errors {
		log.Println("ERROR ", err.Error())
		list = append(list, models.ErrorResponse{
			Code:    ErrorCodes[err].LocalCode,
			Message: err.Error(),
		})
	}
	errorsList := models.ErrorsListResponse{
		Errors: list,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorCodes[errors[0]].HttpCode)
	response, _ := json.Marshal(errorsList)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
