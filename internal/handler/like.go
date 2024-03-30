package handler

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

const USERS_LIKED_LIMIT = 20

func GetPinAndUserId(r *http.Request, ctx context.Context) (entity.PinID, entity.UserID, errs.ErrorInfo) {
	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		return entity.PinID(0), entity.UserID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		}
	}
	pinId := entity.PinID(id)
	userId := ctx.Value("user_id").(entity.UserID)
	//_, userId, err := CheckAuth(r)
	//if err != nil || userId == 0 {
	//	return entity.PinID(0), entity.UserID(0), errs.ErrorInfo{
	//		GeneralErr: err,
	//		LocalErr:   errs.ErrReadCookie,
	//	}
	//}
	return pinId, userId, emptyErrorInfo
}

func (handler *APIHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	l := handler.logger
	ctx := r.Context()

	pinId, userId, errInfo := GetPinAndUserId(r, ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	errInfo = handler.service.SetLike(ctx, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	WriteDefaultResponse(w, l, nil)
}

func (handler *APIHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	l := handler.logger
	ctx := r.Context()

	pinId, userId, errInfo := GetPinAndUserId(r, ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	errInfo = handler.service.ClearLike(ctx, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	WriteDefaultResponse(w, l, nil)
}

func (handler *APIHandler) UsersLiked(w http.ResponseWriter, r *http.Request) {
	l := handler.logger
	ctx := r.Context()

	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pinId := entity.PinID(id)
	res, errInfo := handler.service.GetUsersLiked(ctx, pinId, USERS_LIKED_LIMIT)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}
	WriteDefaultResponse(w, l, res)
}
