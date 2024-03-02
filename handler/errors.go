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
	ErrCheckSession         = errors.New("checking session")
	ErrAlreadyAuthorized    = errors.New("already authorized")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrReadingRequestBody   = errors.New("can't read request body")
	ErrUnmarshalRequestBody = errors.New("can't unmarshal request body")
	ErrInvalidRegisterInput = errors.New("failed to register because of invalid input format")
	ErrInvalidLoginInput    = errors.New("failed to log in because of incorrect email or password")
	ErrHashingPassword      = errors.New("hashing password")
	ErrUserExist            = errors.New("user already exists")
	ErrUserNotExist         = errors.New("user with this email does not exist")
	ErrWrongPassword        = errors.New("wrong password")
	ErrDBUnique             = errors.New("uniqueness error in DB")
	ErrDBInternal           = errors.New("internal db error")
)

// коды ошибок - супер сомнения (мб где-то не 400)
var HttpStatus = map[error]int{
	ErrCheckSession:         400,
	ErrAlreadyAuthorized:    403, // или 405?
	ErrUnauthorized:         401,
	ErrReadingRequestBody:   400,
	ErrUnmarshalRequestBody: 400,
	ErrInvalidRegisterInput: 400,
	ErrInvalidLoginInput:    400,
	ErrHashingPassword:      400,
	ErrUserExist:            400, // мб тоже 401 (после реги мы же автоматом авторизуем)
	ErrUserNotExist:         401,
	ErrWrongPassword:        401,
	ErrDBUnique:             409, // Conflict норм???
	ErrDBInternal:           500,
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
