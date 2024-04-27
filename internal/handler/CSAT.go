package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) AddRating(w http.ResponseWriter, r *http.Request) {

}

func (h *APIHandler) GetRatings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := ctx.Value("request_id").(string)

	rating, err := h.service.GetRating(ctx)
	if err != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestID, err)
		return
	}
	WriteDefaultResponse(w, h.logger, rating)
}

func (h *APIHandler) CreateRatings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := ctx.Value("request_id").(string)

	ratingParams := r.FormValue("rating")
	rating := entity.Rating{}
	err := json.Unmarshal([]byte(ratingParams), &rating)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestID, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	errInfo := h.service.CreateRating(ctx, rating)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestID, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}
