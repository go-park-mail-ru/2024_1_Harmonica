package handler

import (
	"harmonica/internal/entity/errs"
	"log"
	"net/http"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

func MakeErrorInfo(generalErr error, localErr error) errs.ErrorInfo {
	return errs.ErrorInfo{
		GeneralErr: generalErr,
		LocalErr:   localErr,
	}
}

func WriteErrorResponse(w http.ResponseWriter, logger *zap.Logger, requestId string, errInfo errs.ErrorInfo) {
	generalErrMessage := "no general error"
	if errInfo.GeneralErr != nil {
		generalErrMessage = errInfo.GeneralErr.Error()
	}
	if errInfo.LocalErr == nil {
		errInfo.LocalErr = errs.ErrUndefined
	}

	logger.Error(
		errInfo.LocalErr.Error(),
		zap.String("request_id", requestId),
		zap.Int("local_error_code", errs.ErrorCodes[errInfo.LocalErr].LocalCode),
		zap.String("general_error", generalErrMessage),
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errs.ErrorCodes[errInfo.LocalErr].HttpCode)

	response, _ := easyjson.Marshal(errs.ErrorResponse{
		Code:    errs.ErrorCodes[errInfo.LocalErr].LocalCode,
		Message: errInfo.LocalErr.Error(),
	})
	_, err := w.Write(response)
	if err != nil {
		logger.Error("error writing error-response")
		//log.Println(err)
	}
}

func WriteErrorsListResponse(w http.ResponseWriter, logger *zap.Logger, requestId string, errors ...errs.ErrorInfo) {
	var list []errs.ErrorResponse
	for _, err := range errors {
		list = append(list, errs.ErrorResponse{
			Code:    errs.ErrorCodes[err.LocalErr].LocalCode,
			Message: err.LocalErr.Error(),
		})

		generalErrMessage := "no general error"
		if err.GeneralErr != nil {
			generalErrMessage = err.GeneralErr.Error()
		}

		logger.Error(
			err.LocalErr.Error(),
			zap.String("request_id", requestId),
			zap.Int("local_error_code", errs.ErrorCodes[err.LocalErr].LocalCode),
			zap.String("general_error", generalErrMessage),
		)
	}
	errsList := errs.ErrorsListResponse{
		Errors: list,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errs.ErrorCodes[errors[0].LocalErr].HttpCode) // это выглядит как-то не прикольно

	response, _ := easyjson.Marshal(errsList)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
