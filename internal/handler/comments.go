package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value("request_id").(string)
	userId := r.Context().Value("user_id").(entity.UserID)
	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidSlug,
		})
		return
	}
	var comment entity.CommentRequest
	err = UnmarshalRequest(r, &comment)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrReadingRequestBody,
		})
		return
	}
	errInfo := h.service.AddComment(r.Context(), comment.Value, entity.PinID(pinId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
}

func (h *APIHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value("request_id").(string)
	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidSlug,
		})
		return
	}
	res, errInfo := h.service.GetComments(r.Context(), entity.PinID(pinId))
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}
