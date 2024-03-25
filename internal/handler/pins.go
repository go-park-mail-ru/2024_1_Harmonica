package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"log"
	"net/http"
	"strconv"
)

const Limit = 10

func PageToLimitAndOffset(page int) (int, int) {
	return Limit, page * Limit
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

func WriteDefaultResponse(w http.ResponseWriter, object any) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(object)
	_, err := w.Write(response)
	if err != nil {
		log.Print(err)
	}
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
func (handler *APIHandler) GetFeedPins(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}
	pins, err := handler.service.GetFeedPins(r.Context(), limit, offset)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	WriteDefaultResponse(w, pins)
}

func (handler *APIHandler) GetUserPins(w http.ResponseWriter, r *http.Request) {
	userId, err := ReadInt64Slug(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrInvalidSlug)
		return
	}
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}
	pin, err := handler.service.GetUserPins(r.Context(), entity.UserID(userId), limit, offset)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	WriteDefaultResponse(w, pin)
}

func (handler *APIHandler) GetPin(w http.ResponseWriter, r *http.Request) {
	pinId, err := ReadInt64Slug(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrInvalidSlug)
		return
	}
	pin, err := handler.service.GetPinById(r.Context(), entity.PinID(pinId))
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	WriteDefaultResponse(w, pin)
}

func (handler *APIHandler) CreatePin(w http.ResponseWriter, r *http.Request) {
	pin := entity.Pin{AllowComments: true}
	err := UnmarshalRequest(r, &pin)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}
	_, userId, err := CheckAuth(r)
	if err != nil || userId == 0 {
		WriteErrorResponse(w, errs.ErrReadCookie)
		return
	}
	pin.AuthorId = userId
	if pin.ContentUrl == "" {
		WriteErrorResponse(w, errs.ErrInvalidInputFormat) // Тут лучше какую-то другую ошибку
		return
	}
	res, err := handler.service.CreatePin(r.Context(), pin)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	WriteDefaultResponse(w, res)
}

func (handler *APIHandler) UpdatePin(w http.ResponseWriter, r *http.Request) {
	pinId, err := ReadInt64Slug(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrInvalidSlug)
		return
	}
	pin := entity.Pin{AllowComments: true}
	pin.PinId = entity.PinID(pinId)
	err = UnmarshalRequest(r, &pin)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}
	_, userId, err := CheckAuth(r)
	if err != nil || userId == 0 {
		WriteErrorResponse(w, errs.ErrReadCookie)
		return
	}
	pin.AuthorId = userId
	res, err := handler.service.UpdatePin(r.Context(), pin)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	WriteDefaultResponse(w, res)
}

func (handler *APIHandler) DeletePin(w http.ResponseWriter, r *http.Request) {
	pinId, err := ReadInt64Slug(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrInvalidSlug)
		return
	}
	pin := entity.Pin{}
	pin.PinId = entity.PinID(pinId)
	_, userId, err := CheckAuth(r)
	if err != nil || userId == 0 {
		WriteErrorResponse(w, errs.ErrReadCookie)
		return
	}
	pin.AuthorId = userId
	err = handler.service.DeletePin(r.Context(), pin)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	WriteDefaultResponse(w, nil)
}
