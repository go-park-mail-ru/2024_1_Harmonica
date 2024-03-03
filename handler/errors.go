package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrReadCookie         = errors.New("error reading cookie")
	ErrAlreadyAuthorized  = errors.New("already authorized")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrReadingRequestBody = errors.New("error reading request body")
	ErrInvalidInputFormat = errors.New("failed to login/register because of invalid input format")
	ErrHashingPassword    = errors.New("error hashing password")
	ErrUserExist          = errors.New("user with this email already exists")
	ErrUserNotExist       = errors.New("user with this email does not exist")
	ErrWrongPassword      = errors.New("wrong password")
	ErrDBUniqueEmail      = errors.New("user with this email already exists")
	ErrDBUniqueNickname   = errors.New("user with this nickname already exists")
	ErrDBInternal         = errors.New("internal db error")
)

var HttpStatus = map[error]int{
	ErrReadCookie:         400,
	ErrAlreadyAuthorized:  403,
	ErrUnauthorized:       401,
	ErrReadingRequestBody: 400,
	ErrInvalidInputFormat: 400,
	ErrHashingPassword:    400,
	ErrUserExist:          400,
	ErrUserNotExist:       401,
	ErrWrongPassword:      401,
	ErrDBUniqueEmail:      500,
	ErrDBUniqueNickname:   500,
	ErrDBInternal:         500,
}

type errorResponse struct {
	Message string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	log.Print("ERROR ", err.Error())
	w.WriteHeader(HttpStatus[err])
	response, _ := json.Marshal(errorResponse{Message: err.Error()})
	w.Write(response) // Unhandled error - наверное тут без разницы
}
