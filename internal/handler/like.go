package handler

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

const USERS_LIKED_LIMIT = 20

func GetPinAndUserId(r *http.Request, ctx context.Context) (entity.PinID, entity.UserID, errs.ErrorInfo) {
	id, err := ReadUint64Slug(r, "pin_id")
	if err != nil {
		return entity.PinID(0), entity.UserID(0), errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		}
	}
	pinId := entity.PinID(id)
	userId := ctx.Value("user_id").(entity.UserID)
	return pinId, userId, emptyErrorInfo
}

func (h *APIHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pinId, userId, errInfo := GetPinAndUserId(r, ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	errInfo = h.service.SetLike(ctx, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pinId, userId, errInfo := GetPinAndUserId(r, ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	errInfo = h.service.ClearLike(ctx, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

func (h *APIHandler) UsersLiked(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := ReadUint64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pinId := entity.PinID(id)
	res, errInfo := h.service.GetUsersLiked(ctx, pinId, USERS_LIKED_LIMIT)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}
