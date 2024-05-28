package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
	"time"
)

func (h *APIHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("request_id").(string)
	if !ok {
		requestId = "0"
	}
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
	pin, commentId, errInfo := h.service.AddComment(r.Context(), comment.Value, entity.PinID(pinId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}

	// создание уведомления о новом комментарии для автора пина
	n := entity.Notification{
		Type:              entity.NotificationTypeComment,
		UserId:            pin.PinAuthor.UserId,
		TriggeredByUserId: userId,
		CommentId:         commentId,
		PinId:             entity.PinID(pinId),
	}
	nId, errInfo := h.service.CreateNotification(r.Context(), n)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}

	// инфа о юзере, который прокомментировал
	user, errInfo := h.service.GetUserById(r.Context(), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	// отправка в websocket - отправляем автору пина
	notification := &entity.WSMessage{
		Action: entity.WSActionNotificationComment,
		Payload: entity.WSMessagePayload{
			UserId: pin.PinAuthor.UserId,
			TriggeredByUser: entity.TriggeredByUser{
				UserId:    userId,
				Nickname:  user.Nickname,
				AvatarURL: user.AvatarURL,
			},
			NotificationId: nId,
			Comment: entity.CommentNotificationResponse{
				Text: comment.Value,
			},
			Pin: entity.PinNotificationResponse{
				PinId:      entity.PinID(pinId),
				ContentUrl: pin.ContentUrl,
			},
			CreatedAt: time.Now(),
		},
	}
	h.hub.broadcast <- notification

	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	requestId, ok := r.Context().Value("request_id").(string)
	if !ok {
		requestId = "0"
	}
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
