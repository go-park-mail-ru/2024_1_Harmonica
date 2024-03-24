package handler

import (
	"encoding/json"
	"harmonica/internal/entity/errs"
	"log"
	"net/http"
	"strconv"
)

const Limit = 10

func PageToLimitAndOffset(page int) (int, int) {
	return Limit, page * Limit
}

// Pins List
//
//	@Summary		Pins list
//	@Description	Get pins by page
//	@Tags			Pins
//	@Param			page	query		int	false	"Page num from 0"
//	@Success		200		{object}	entity.Pins
//	@Failure		400		{object}	errs.ErrorResponse
//	@Router			/pins_list [get]
func (handler *APIHandler) PinsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		pageString = "0"
	}
	page, err := strconv.Atoi(pageString)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}
	limit, offset := PageToLimitAndOffset(page)
	pins, err := handler.service.GetPins(ctx, limit, offset)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(pins)
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
