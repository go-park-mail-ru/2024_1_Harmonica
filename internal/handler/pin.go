package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

const FEED_PINS_LIMIT = 40

func PageToLimitAndOffset(page int) (int, int) {
	return FEED_PINS_LIMIT, page * FEED_PINS_LIMIT
}

func GetLimitAndOffset(r *http.Request) (int, int, error) {
	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		pageString = "0"
	}
	page, err := strconv.Atoi(pageString)
	if err != nil {
		return 0, 0, err
	}
	limit, offset := PageToLimitAndOffset(page)
	return limit, offset, nil
}

func UnmarshalRequest(r *http.Request, dest easyjson.Unmarshaler) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = easyjson.Unmarshal(bodyBytes, dest)
	return err
}

func WriteDefaultResponse(w http.ResponseWriter, logger *zap.Logger, object easyjson.Marshaler) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := easyjson.Marshal(object)
	_, err := w.Write(response)
	if err != nil {
		logger.Error(
			errs.ErrServerInternal.Error(),
			zap.Int("local_error_code", errs.ErrorCodes[errs.ErrServerInternal].LocalCode),
			zap.String("general_error", err.Error()),
		)
	}
}

// Feed pins list
//
//	@Summary		Pins list
//	@Description	Get pins by page
//	@Tags			Pins
//	@Param			page	query		int	false	"Page num from 0"
//	@Success		200		{object}	entity.FeedPins
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 4."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins [get]
func (h *APIHandler) Feed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	feedType := r.URL.Query().Get("type")
	var (
		pins    entity.FeedPins
		errInfo errs.ErrorInfo
		userId  entity.UserID
		ok      bool
	)
	switch feedType {
	case "subscriptions":
		userIdFromSession := ctx.Value("user_id")
		if userIdFromSession != nil {
			userId, ok = userIdFromSession.(entity.UserID)
			if !ok {
				userId = 0
			}
		}
		if userId != 0 {
			pins, errInfo = h.service.GetSubscriptionsFeedPins(ctx, userId, limit, offset)
			break
		}
		//pins, errInfo = h.service.GetFeedPins(ctx, limit, offset)
		fallthrough
	default:
		pins, errInfo = h.service.GetFeedPins(ctx, limit, offset)
	}

	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, pins)
}

// Get pins that created by user id.
//
//	@Summary		Get pins that created by user id
//	@Description	Get pins of user by page
//	@Tags			Pins
//	@Param			page	query		int	false	"Page num from 0"
//	@Success		200		{object}	entity.FeedPins
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 4, 12."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins/created/{nickname} [get]
func (h *APIHandler) UserPins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	userNickname, err := ReadStringSlug(r, "nickname")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	pin, errInfo := h.service.GetUserPins(r.Context(), userNickname, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, pin)
}

// Get pin by id.
//
//	@Summary		Get pin by id
//	@Description	Get pin by id in the slug
//	@Tags			Pins
//	@Param			pin_id	path		int	true	"Pin ID"
//	@Success		200		{object}	entity.PinPageResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 4, 12."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins/{pin_id} [get]
func (h *APIHandler) GetPin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	var userId entity.UserID
	userIdFromSession := ctx.Value("user_id")
	if userIdFromSession != nil {
		id, ok := userIdFromSession.(entity.UserID)
		if ok {
			userId = id
		}
	}
	pin, errInfo := h.service.GetPinById(ctx, entity.PinID(pinId), userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	pin.IsOwner = pin.PinAuthor.UserId == userId
	WriteDefaultResponse(w, h.logger, pin)
}

// Create pin.
//
//	@Summary		Create pin
//	@Description	Create pin by description
//	@Tags			Pins
//	@Produce		json
//	@Accept			multipart/form-data
//	@Param			Cookie	header		string	true	"session-token"	default(session-token=)
//	@Param			pin		formData	string	false	"Pin information in json"
//	@Param			image	formData	file	true	"Pin image"
//	@Success		200		{object}	entity.PinPageResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 15, 18, 19."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins [post]
func (h *APIHandler) CreatePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	pinParams := r.FormValue("pin")
	pin := entity.Pin{AllowComments: true}
	err := easyjson.Unmarshal([]byte(pinParams), &pin)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	pin.AuthorId = ctx.Value("user_id").(entity.UserID)
	imageName, err := h.UploadImage(r, "image")
	if err != nil {
		localErr := err
		if errs.ErrorCodes[localErr].HttpCode == 0 {
			localErr = errs.ErrDBInternal
		}
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{GeneralErr: err, LocalErr: localErr})
		return
	}
	pin.ContentUrl = h.FormImgURL(imageName)

	res, errInfo := h.service.CreatePin(ctx, pin)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}

	// создание уведомления для подписчиках юзера, выложившего пин
	n := entity.Notification{
		Type:              entity.NotificationTypeNewPin,
		TriggeredByUserId: pin.AuthorId,
		PinId:             res.PinId,
	}
	_, errInfo = h.service.CreateNotification(ctx, n)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}

	//// подписчики - это те, кому нужно отправить уведомление о публикации пина
	//subscribers, errInfo := h.service.GetUserSubscribers(ctx, pin.AuthorId)
	//if errInfo != emptyErrorInfo {
	//	WriteErrorResponse(w, h.logger, requestId, errInfo)
	//	return
	//}
	//// отправка в websocket
	//notification := &entity.WSMessage{
	//	Action: entity.WSActionNotificationNewPin,
	//	Payload: entity.WSMessagePayload{
	//		TriggeredByUser: entity.TriggeredByUser{
	//			UserId:    res.PinAuthor.UserId,
	//			Nickname:  res.PinAuthor.Nickname,
	//			AvatarURL: res.PinAuthor.AvatarURL,
	//		},
	//		NotificationId: nId,
	//		Pin: entity.PinNotificationResponse{
	//			PinId:      res.PinId,
	//			ContentUrl: res.ContentUrl,
	//		},
	//		CreatedAt: time.Now(),
	//	},
	//}
	//for _, s := range subscribers.Subscribers {
	//	notification.Payload.UserId = s.UserId
	//	h.hub.broadcast <- notification
	//}

	WriteDefaultResponse(w, h.logger, res)
}

// Update pin.
//
//	@Summary		Update pin
//	@Description	Update pin by description
//	@Tags			Pins
//	@Produce		json
//	@Accept			json
//	@Param			pin		body		entity.Pin	true	"Pin information"
//	@Param			Cookie	header		string		true	"session-token"	default(session-token=)
//	@Success		200		{object}	entity.PinPageResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 12"
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins/{pin_id} [post]
func (h *APIHandler) UpdatePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pin := entity.Pin{AllowComments: true}
	pin.PinId = entity.PinID(pinId)
	err = UnmarshalRequest(r, &pin)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	pin.AuthorId = ctx.Value("user_id").(entity.UserID)
	res, errInfo := h.service.UpdatePin(ctx, pin)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}

// Delete pin.
//
//	@Summary		Delete pin
//	@Description	Delete pin by id (allowed only for pin's author)
//	@Tags			Pins
//	@Param			Cookie	header		string	true	"session-token"	default(session-token=)
//	@Param			pin_id	path		int		true	"Pin ID"
//	@Success		200		{object}	entity.PinPageResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4 12"
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins/{pin_id} [delete]
func (h *APIHandler) DeletePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pin := entity.Pin{}
	pin.PinId = entity.PinID(pinId)
	pin.AuthorId = ctx.Value("user_id").(entity.UserID)
	errInfo := h.service.DeletePin(ctx, pin)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}
