package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errors_list"
	"log"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, err error) {
	log.Println("ERROR ", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errors_list.ErrorCodes[err].HttpCode)
	response, _ := json.Marshal(entity.ErrorResponse{
		Code:    errors_list.ErrorCodes[err].LocalCode,
		Message: err.Error(),
	})
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
}

func WriteErrorsListResponse(w http.ResponseWriter, errors ...error) {
	var list []entity.ErrorResponse
	for _, err := range errors {
		log.Println("ERROR ", err.Error())
		list = append(list, entity.ErrorResponse{
			Code:    errors_list.ErrorCodes[err].LocalCode,
			Message: err.Error(),
		})
	}
	eList := entity.ErrorsListResponse{
		Errors: list,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errors_list.ErrorCodes[errors[0]].HttpCode)
	response, _ := json.Marshal(eList)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
