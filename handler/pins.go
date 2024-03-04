package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const Limit = 10

func PageToLimitAndOffset(page int) (int, int) {
	return Limit, page * Limit
}

func (handler *APIHandler) PinsList(w http.ResponseWriter, r *http.Request) {
	if curSessionToken, err := CheckAuth(r); err != nil || curSessionToken == "" {
		WriteErrorResponse(w, ErrUnauthorized)
		return
	}
	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		pageString = "0"
	}
	page, err := strconv.Atoi(pageString)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}
	limit, offset := PageToLimitAndOffset(page)
	pins, err := handler.connector.GetPins(limit, offset)
	if err != nil {
		WriteErrorResponse(w, ErrDBInternal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(pins)
	w.Write(response)
}
