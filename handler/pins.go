package handler

import (
	"encoding/json"
	"harmonica/db"
	"log"
	"net/http"
	"strconv"
)

const Limit = 10

func LogIfError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func PageToLimitAndOffset(page int) (int, int) {
	return Limit, page * Limit
}

func (handler *APIHandler) PinsList(w http.ResponseWriter, r *http.Request) {

	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		pageString = "0"
	}
	page, err := strconv.Atoi(pageString)
	if err != nil {
		log.Print(err)
		SetHttpError(w, ErrReadingRequestBody, 400)
		return
	}

	limit, offset := PageToLimitAndOffset(page)
	pins, err := handler.connector.GetPins(limit, offset)
	if err != nil {
		log.Print(err)
		SetHttpError(w, ErrDB, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		Pins []db.Pin `json:"pins"`
	}{Pins: pins})
	w.Write(response)
}
