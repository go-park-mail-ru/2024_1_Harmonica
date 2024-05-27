package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
	"time"
)

func (h *APIHandler) SubscribeToUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)
	userId, subscribeUserId, errInfo := GetSubscriptionInfoFromSlugAndContext(r)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	errInfo = h.service.AddSubscriptionToUser(ctx, userId, subscribeUserId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}

	// инфа о юзере, который подписывается
	user, errInfo := h.service.GetUserById(ctx, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	// отправка в websocket - отправляем тому, на кого подписались
	notification := &entity.WSMessage{
		Action: entity.WSActionNotificationSubscription,
		Payload: entity.WSMessagePayload{
			UserId: subscribeUserId,
			TriggeredByUser: entity.TriggeredByUser{
				UserId:    userId,
				Nickname:  user.Nickname,
				AvatarURL: user.AvatarURL,
			},
			CreatedAt: time.Now(),
		},
	}
	h.hub.broadcast <- notification

	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) UnsubscribeFromUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)
	userId, unsubscribeUserId, errInfo := GetSubscriptionInfoFromSlugAndContext(r)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	errInfo = h.service.DeleteSubscriptionToUser(ctx, userId, unsubscribeUserId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) GetUserSubscribers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)
	userId, err := ReadInt64Slug(r, "user_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	subscribers, errInfo := h.service.GetUserSubscribers(ctx, entity.UserID(userId))
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, subscribers)
}

func (h *APIHandler) GetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)
	userId, err := ReadInt64Slug(r, "user_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}
	subscribers, errInfo := h.service.GetUserSubscriptions(ctx, entity.UserID(userId))
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, subscribers)
}

func GetSubscriptionInfoFromSlugAndContext(r *http.Request) (entity.UserID, entity.UserID, errs.ErrorInfo) {
	userIdFromSession, ok := r.Context().Value("user_id").(entity.UserID)
	if !ok {
		return 0, 0, MakeErrorInfo(nil, errs.ErrTypeConversion)
	}
	userIdFromSlug, err := ReadInt64Slug(r, "user_id")
	if err != nil {
		return 0, 0, MakeErrorInfo(err, errs.ErrInvalidSlug)
	}
	return userIdFromSession, entity.UserID(userIdFromSlug), errs.ErrorInfo{}
}
