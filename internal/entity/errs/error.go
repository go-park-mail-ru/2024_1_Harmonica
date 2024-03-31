package errs

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorsListResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorInfo struct {
	GeneralErr error
	LocalErr   error
}
