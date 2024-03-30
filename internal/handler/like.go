package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

const USERS_LIKED_LIMIT = 20

var emptyErrorInfo = errs.ErrorInfo{}

func GetPinAndUserId(r *http.Request) (entity.PinID, entity.UserID, errs.ErrorInfo) {
	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		return entity.PinID(0), entity.UserID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		}
	}
	pinId := entity.PinID(id)
	_, userId, err := CheckAuth(r)
	if err != nil || userId == 0 {
		return entity.PinID(0), entity.UserID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadCookie,
		}
	}
	return pinId, userId, emptyErrorInfo
}

func (handler *APIHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	l := handler.logger

	pinId, userId, errInfo := GetPinAndUserId(r)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	errInfo = handler.service.SetLike(r.Context(), pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	WriteDefaultResponse(w, nil)
}

func (handler *APIHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	l := handler.logger
	pinId, userId, errInfo := GetPinAndUserId(r)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	errInfo = handler.service.ClearLike(r.Context(), pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	WriteDefaultResponse(w, nil)
}

func (handler *APIHandler) UsersLiked(w http.ResponseWriter, r *http.Request) {
	l := handler.logger

	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pinId := entity.PinID(id)
	res, errInfo := handler.service.GetUsersLiked(r.Context(), pinId, USERS_LIKED_LIMIT)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	WriteDefaultResponse(w, res)
}
