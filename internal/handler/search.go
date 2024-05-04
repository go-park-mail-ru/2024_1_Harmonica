package handler

import (
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	query, err := ReadStringSlug(r, "search_query")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: err})
		return
	}
	res, errInfo := h.service.Search(ctx, query)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}
