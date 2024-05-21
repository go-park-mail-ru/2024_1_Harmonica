package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, ok := r.Context().Value("request_id").(string)
	if !ok {
		requestId = "0"
	}

	message := entity.Message{}
	receiverId, err := ReadInt64Slug(r, "receiver_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	message.ReceiverId = entity.UserID(receiverId)

	err = UnmarshalRequest(r, &message)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	if !ValidateMessage(message) {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}

	senderId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
	}
	message.SenderId = senderId

	errInfo := h.service.CreateMessage(ctx, message)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) ReadMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, ok := r.Context().Value("request_id").(string)
	if !ok {
		requestId = "0"
	}

	dialogUserId, err := ReadInt64Slug(r, "user_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	authUserId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
	}

	messages, errInfo := h.service.GetMessages(ctx, entity.UserID(dialogUserId), authUserId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, messages)
}

func (h *APIHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
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
	chats, errInfo := h.service.GetUserChats(ctx, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, chats)
}
