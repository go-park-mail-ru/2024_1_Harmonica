package errs

import "errors"

var (
	ErrAlreadyAuthorized   = errors.New("already authorized")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrReadCookie          = errors.New("error reading cookie")
	ErrReadingRequestBody  = errors.New("error reading request body")
	ErrInvalidInputFormat  = errors.New("validation conditions are not met")
	ErrHashingPassword     = errors.New("error hashing password")
	ErrUserNotExist        = errors.New("user does not exist")
	ErrWrongPassword       = errors.New("wrong password (can't authorize)")
	ErrDBUniqueEmail       = errors.New("user with this email already exists (can't register)")
	ErrDBUniqueNickname    = errors.New("user with this nickname already exists (can't register)")
	ErrDBInternal          = errors.New("internal db error")
	ErrInvalidSlug         = errors.New("invalid slug parameter")
	ErrDiffUserId          = errors.New("user id in slug and session don't match")
	ErrPermissionDenied    = errors.New("current user doesn't have requested permissions")
	ErrContentUrlRequired  = errors.New("field 'content_url' is required")
	ErrEmptyContentURL     = errors.New("content url can't be empty")
	ErrInvalidContentType  = errors.New("invalid content type header")
	ErrInvalidImg          = errors.New("image is not valid")
	ErrNoImageProvided     = errors.New("there is not any image file")
	ErrElementNotExist     = errors.New("element does not exist")
	ErrTypeConversion      = errors.New("type conversion error")
	ErrDBUniqueViolation   = errors.New("unique violation error (element already exists in db)")
	ErrGRPCWentWrong       = errors.New("something went wrong in grpc connection")
	ErrCantParseTime       = errors.New("cant parse timestamp")
	ErrForeignKeyViolation = errors.New("foreign key violation error (referenced element does not exist in db)")
	ErrWSConnectionClosed  = errors.New("ws connection was closed")
	ErrWSConnectionUpgrade = errors.New("can't upgrade http connection to ws")
	ErrServerInternal      = errors.New("internal server error")
	ErrNotFound            = errors.New("not found")
	ErrMinioTurnedOff      = errors.New("minio server doesn't response")
	ErrNotAllowedExtension = errors.New("this file extension is not allowed")
	ErrUndefined           = errors.New("this error is undefined")
	ErrEmptyComment        = errors.New("empty comment")
	// для случая, если забудем передать значение LocalErr в errInfo;
	// до этого нововведения при таком раскладе событий приложение паниковало
)

var ErrorCodes = map[error]struct {
	HttpCode  int
	LocalCode int
}{
	ErrAlreadyAuthorized:   {HttpCode: 403, LocalCode: 1},
	ErrUnauthorized:        {HttpCode: 401, LocalCode: 2},
	ErrReadCookie:          {HttpCode: 400, LocalCode: 3},
	ErrReadingRequestBody:  {HttpCode: 400, LocalCode: 4},
	ErrInvalidInputFormat:  {HttpCode: 400, LocalCode: 5},
	ErrHashingPassword:     {HttpCode: 500, LocalCode: 6},
	ErrUserNotExist:        {HttpCode: 404, LocalCode: 7},
	ErrWrongPassword:       {HttpCode: 401, LocalCode: 8},
	ErrDBUniqueEmail:       {HttpCode: 500, LocalCode: 9},
	ErrDBUniqueNickname:    {HttpCode: 500, LocalCode: 10},
	ErrDBInternal:          {HttpCode: 500, LocalCode: 11},
	ErrInvalidSlug:         {HttpCode: 400, LocalCode: 12},
	ErrDiffUserId:          {HttpCode: 400, LocalCode: 13},
	ErrPermissionDenied:    {HttpCode: 403, LocalCode: 14},
	ErrContentUrlRequired:  {HttpCode: 400, LocalCode: 15},
	ErrEmptyContentURL:     {HttpCode: 400, LocalCode: 16},
	ErrInvalidContentType:  {HttpCode: 400, LocalCode: 17},
	ErrInvalidImg:          {HttpCode: 400, LocalCode: 18},
	ErrNoImageProvided:     {HttpCode: 400, LocalCode: 19},
	ErrElementNotExist:     {HttpCode: 400, LocalCode: 20},
	ErrTypeConversion:      {HttpCode: 400, LocalCode: 21},
	ErrDBUniqueViolation:   {HttpCode: 500, LocalCode: 22},
	ErrGRPCWentWrong:       {HttpCode: 500, LocalCode: 23},
	ErrCantParseTime:       {HttpCode: 500, LocalCode: 24},
	ErrForeignKeyViolation: {HttpCode: 500, LocalCode: 25},
	ErrWSConnectionClosed:  {HttpCode: 500, LocalCode: 26},
	ErrWSConnectionUpgrade: {HttpCode: 500, LocalCode: 27},
	ErrServerInternal:      {HttpCode: 500, LocalCode: 28},
	ErrNotFound:            {HttpCode: 404, LocalCode: 29},
	ErrMinioTurnedOff:      {HttpCode: 500, LocalCode: 30},
	ErrNotAllowedExtension: {HttpCode: 400, LocalCode: 31},
	ErrUndefined:           {HttpCode: 400, LocalCode: 32},
	ErrEmptyComment:        {HttpCode: 400, LocalCode: 33},
}

var GetLocalErrorByCode = map[int64]error{
	1:  ErrAlreadyAuthorized,
	2:  ErrUnauthorized,
	3:  ErrReadCookie,
	4:  ErrReadingRequestBody,
	5:  ErrInvalidInputFormat,
	6:  ErrHashingPassword,
	7:  ErrUserNotExist,
	8:  ErrWrongPassword,
	9:  ErrDBUniqueEmail,
	10: ErrDBUniqueNickname,
	11: ErrDBInternal,
	12: ErrInvalidSlug,
	13: ErrDiffUserId,
	14: ErrPermissionDenied,
	15: ErrContentUrlRequired,
	16: ErrEmptyContentURL,
	17: ErrInvalidContentType,
	18: ErrInvalidImg,
	19: ErrNoImageProvided,
	20: ErrElementNotExist,
	21: ErrTypeConversion,
	22: ErrDBUniqueViolation,
	23: ErrGRPCWentWrong,
	24: ErrCantParseTime,
	25: ErrForeignKeyViolation,
	26: ErrWSConnectionClosed,
	27: ErrWSConnectionUpgrade,
	28: ErrServerInternal,
	29: ErrNotFound,
	30: ErrMinioTurnedOff,
	31: ErrNotAllowedExtension,
	32: ErrUndefined,
	33: ErrEmptyComment,
}
