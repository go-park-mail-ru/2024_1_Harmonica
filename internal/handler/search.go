package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"
)

func (h *APIHandler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr:   errs.ErrReadingRequestBody,
			GeneralErr: err,
		})
		return
	}
	var req entity.SearchRequest
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr:   errs.ErrReadingRequestBody,
			GeneralErr: err,
		})
		return
	}

	res, errInfo := h.service.Search(ctx, req)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}
