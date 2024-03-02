package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	ErrReadingRequestBody   = errors.New("error reading request body")
	ErrUnmarshalRequestBody = errors.New("error unmarshal request body")
	ErrInvalidRegisterInput = errors.New("failed to register because of invalid input format")
	ErrInvalidLoginInput    = errors.New("failed to log in because of incorrect email or password")
	ErrHashingPassword      = errors.New("error hashing password")
	ErrUserExist            = errors.New("user already exists")
	ErrDB                   = errors.New("db error")
)

var HttpStatusByErr = map[error]int{
	ErrReadingRequestBody: 400,
	ErrDB:                 500,
}

func SetHttpError(w http.ResponseWriter, serverErr error, localErr error, status int) {
	fmtMessage := ""
	switch {
	case localErr == nil:
		fmtMessage = fmt.Sprintf("ERROR %s", serverErr.Error())
	default:
		fmtMessage = fmt.Sprintf("ERROR %s: %s", serverErr.Error(), localErr.Error())
	}
	log.Print(fmtMessage)
	http.Error(w, fmtMessage, status)
}
