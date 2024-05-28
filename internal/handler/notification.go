package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) GetUnreadNotifications(w http.ResponseWriter, r *http.Request) {
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

func (h *APIHandler) ReadNotification(w http.ResponseWriter, r *http.Request) {
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
	notificationId, err := ReadInt64Slug(r, "notification_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	errInfo := h.service.ReadNotification(ctx, entity.NotificationID(notificationId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) ReadAllNotifications(w http.ResponseWriter, r *http.Request) {
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
	errInfo := h.service.ReadAllNotifications(ctx, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}
