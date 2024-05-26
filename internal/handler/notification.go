package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) ReadNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, ok := r.Context().Value("request_id").(string)
	if !ok {
		requestId = "0"
	}
	userId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
	}

	notifications, errInfo := h.service.GetUnreadNotifications(ctx, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, notifications)
}
