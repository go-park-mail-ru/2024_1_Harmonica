package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

//type HttpError struct {
//	err    error
//	status int
//}

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
	ErrReadingRequestBody:   400,
	ErrUnmarshalRequestBody: 400,
	ErrInvalidRegisterInput: 400,
	ErrInvalidLoginInput:    400,
	ErrHashingPassword:      400,
	ErrUserExist:            400,
	ErrDB:                   400,
}

func SetHttpError(w http.ResponseWriter, serverErr error, localErr error, status int) {
	fmtMessage := ""
	switch {
	case localErr == nil:
		fmtMessage = fmt.Sprintf("ERROR %s", serverErr.Error())
	default:
		fmtMessage = fmt.Sprintf("ERROR %s: %s", serverErr.Error(), localErr.Error())
	}
	log.Println(fmtMessage)
	http.Error(w, fmtMessage, status)
}
