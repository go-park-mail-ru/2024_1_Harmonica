package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func (h *APIHandler) UpdateDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, ok := r.Context().Value("request_id").(string)
	if !ok {
		requestId = "0"
	}

	draft := entity.Draft{}
	receiverId, err := ReadInt64Slug(r, "receiver_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	draft.ReceiverId = entity.UserID(receiverId)
	err = UnmarshalRequest(r, &draft)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	// решили без валидации
	senderId, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrTypeConversion))
		return
	}
	draft.SenderId = senderId

	errInfo := h.service.UpdateDraft(ctx, draft)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}
