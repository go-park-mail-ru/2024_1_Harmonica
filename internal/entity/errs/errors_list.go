package errs

import "errors"

var (
	ErrAlreadyAuthorized  = errors.New("already authorized")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrReadCookie         = errors.New("error reading cookie")
	ErrReadingRequestBody = errors.New("error reading request body")
	ErrInvalidInputFormat = errors.New("validation conditions are not met")
	ErrHashingPassword    = errors.New("error hashing password")
	ErrUserNotExist       = errors.New("user with this email does not exist (can't authorize)")
	ErrWrongPassword      = errors.New("wrong password (can't authorize)")
	ErrDBUniqueEmail      = errors.New("user with this email already exists (can't register)")
	ErrDBUniqueNickname   = errors.New("user with this nickname already exists (can't register)")
	ErrDBInternal         = errors.New("internal db error")
	ErrInvalidSlug        = errors.New("invalid slug parameter")                  // NEW !
	ErrDiffUserId         = errors.New("user id in slug and session don't match") // NEW !
	ErrPermissionDenied   = errors.New("current user doesn't have requested permissions")
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
	ErrInvalidSlug:        {HttpCode: 400, LocalCode: 12},
	ErrDiffUserId:         {HttpCode: 400, LocalCode: 13},
	ErrPermissionDenied:   {HttpCode: 403, LocalCode: 14},
}
