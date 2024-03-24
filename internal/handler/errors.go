package handler

import (
	"encoding/json"
	"harmonica/internal/entity/errs"
	"log"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, err error) {
	log.Println("ERROR", err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errs.ErrorCodes[err].HttpCode)

	response, _ := json.Marshal(errs.ErrorResponse{
		Code:    errs.ErrorCodes[err].LocalCode,
		Message: err.Error(),
	})
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
}

func WriteErrorsListResponse(w http.ResponseWriter, errors ...error) {
	var list []errs.ErrorResponse
	for _, err := range errors {
		log.Println("ERROR", err.Error())
		list = append(list, errs.ErrorResponse{
			Code:    errs.ErrorCodes[err].LocalCode,
			Message: err.Error(),
		})
	}
	errsList := errs.ErrorsListResponse{
		Errors: list,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errs.ErrorCodes[errors[0]].HttpCode)

	response, _ := json.Marshal(errsList)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
