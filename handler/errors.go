package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type HttpError struct {
	err    error
	status int
}

var (
	ErrReadingRequestBody   = errors.New("error reading request body")
	ErrUnmarshalRequestBody = errors.New("error unmarshal request body")
	ErrInvalidRegisterInput = errors.New("failed to register because of invalid input format")
	ErrInvalidLogin         = errors.New("failed to log in because of incorrect email or password")
	ErrHashingPassword      = errors.New("error hashing password")
	ErrUserExist            = errors.New("user already exists")
	ErrDB                   = errors.New("db error")
)

var HttpErrSlice = []HttpError{
	{
		err:    ErrReadingRequestBody,
		status: 400,
	},
	{
		err:    ErrUnmarshalRequestBody,
		status: 400,
	},
	{
		err:    ErrInvalidRegisterInput,
		status: 400,
	},
	{
		err:    ErrHashingPassword,
		status: 400,
	},
	{
		err:    ErrUserExist,
		status: 400,
	},
	{
		err:    ErrDB,
		status: 400,
	},
}

func SetHttpError(w http.ResponseWriter, httpError HttpError, err error) {
	if err == nil {
		log.Println(httpError.err)
		http.Error(w, httpError.err.Error(), httpError.status)
		return
	}
	log.Printf("%s: %s", httpError.err, err.Error())
	fmtMessage := fmt.Sprintf("%s: %s", httpError.err.Error(), err.Error())
	http.Error(w, fmtMessage, httpError.status)
}
