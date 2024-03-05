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

// Pins List
//	@Summary		Pins list
//	@Description	Get pins by page
//	@Tags			Pins
//	@Param			page	query		int	false	"Page num from 0"
//	@Success		200		{object}	models.Pins
//	@Failure		400		{object}	errorResponse
//	@Router			/pins_list [get]
func (handler *APIHandler) PinsList(w http.ResponseWriter, r *http.Request) {
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
