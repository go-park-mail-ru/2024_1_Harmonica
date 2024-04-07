package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

const FEED_PINS_LIMIT = 10

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

func UnmarshalRequest(r *http.Request, dest any) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &dest)
	return err
}

func WriteDefaultResponse(w http.ResponseWriter, logger *zap.Logger, object any) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(object)
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
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	pins, errInfo := h.service.GetFeedPins(r.Context(), limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
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
	userNickanme, err := ReadStringSlug(r, "nickname")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	pin, errInfo := h.service.GetUserPins(r.Context(), userNickanme, limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
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
	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pin, errInfo := h.service.GetPinById(ctx, entity.PinID(pinId))
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	if ctx.Value("is_auth") == true {
		userIdFromSession, ok := ctx.Value("user_id").(entity.UserID)
		if !ok {
			WriteErrorResponse(w, h.logger, errs.ErrorInfo{
				LocalErr: errs.ErrTypeConversion,
			})
		}
		pin.IsOwner = pin.PinAuthor.UserId == userIdFromSession
	}
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
	pinParams := r.FormValue("pin")
	pin := entity.Pin{AllowComments: true}
	err := json.Unmarshal([]byte(pinParams), &pin)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
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
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{GeneralErr: err, LocalErr: localErr})
		return
	}
	pin.ContentUrl = FormImgURL(imageName)

	res, errInfo := h.service.CreatePin(ctx, pin)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
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

	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pin := entity.Pin{AllowComments: true}
	pin.PinId = entity.PinID(pinId)
	err = UnmarshalRequest(r, &pin)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	pin.AuthorId = ctx.Value("user_id").(entity.UserID)
	res, errInfo := h.service.UpdatePin(ctx, pin)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
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

	pinId, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
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
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}
