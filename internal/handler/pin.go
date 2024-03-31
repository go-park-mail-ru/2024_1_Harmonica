package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"
	"strconv"
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

// Feed pins list
//
//	@Summary		Pins list
//	@Description	Get pins by page
//	@Tags			Pins
//	@Param			page	query		int	false	"Page num from 0"
//	@Success		200		{object}	entity.Pins
//	@Failure		400		{object}	errs.ErrorResponse
//	@Router			/pins_list [get]
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

func (h *APIHandler) UserPins(w http.ResponseWriter, r *http.Request) {
	userId, err := ReadUint64Slug(r, "user_id")
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
	pin, errInfo := h.service.GetUserPins(r.Context(), entity.UserID(userId), limit, offset)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, pin)
}

func (h *APIHandler) GetPin(w http.ResponseWriter, r *http.Request) {
	pinId, err := ReadUint64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pin, errInfo := h.service.GetPinById(r.Context(), entity.PinID(pinId))
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, emptyErrorInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, pin)
}

func (h *APIHandler) CreatePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pin := entity.Pin{AllowComments: true}
	err := UnmarshalRequest(r, &pin)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	pin.AuthorId = ctx.Value("user_id").(entity.UserID)
	if pin.ContentUrl == "" {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrEmptyContentURL,
		})
		return
	}
	res, errInfo := h.service.CreatePin(ctx, pin)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, emptyErrorInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}

func (h *APIHandler) UpdatePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pinId, err := ReadUint64Slug(r, "pin_id")
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

func (h *APIHandler) DeletePin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pinId, err := ReadUint64Slug(r, "pin_id")
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
