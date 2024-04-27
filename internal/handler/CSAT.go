package handler

import (
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
