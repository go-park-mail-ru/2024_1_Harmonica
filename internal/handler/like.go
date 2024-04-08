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
	return pinId, userId, emptyErrorInfo
}

// Set like on the pin
//
//	@Summary		Set like on the pin
//	@Description	Sets like by pin id and auth token
//	@Tags			Likes
//	@Param			Cookie	header	string	true	"session-token"	default(session-token=)
//	@Param			pin_id	path	int		true	"Pin ID"
//	@Success		200
//	@Failure		400	{object}	errs.ErrorResponse	"Possible code responses: 3, 12."
//	@Failure		401	{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		500	{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins/{pin_id}/like [post]
func (h *APIHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	pinId, userId, errInfo := GetPinAndUserId(r, ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	errInfo = h.service.SetLike(ctx, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

// Delete like on the pin
//
//	@Summary		Delete like on the pin
//	@Description	Delete like by pin id and auth token
//	@Tags			Likes
//	@Param			Cookie	header	string	true	"session-token"	default(session-token=)
//	@Param			pin_id	path	int		true	"Pin ID"
//	@Success		200
//	@Failure		400	{object}	errs.ErrorResponse	"Possible code responses: 3, 12."
//	@Failure		401	{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		500	{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/pins/{pin_id}/like [delete]
func (h *APIHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	pinId, userId, errInfo := GetPinAndUserId(r, ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	errInfo = h.service.ClearLike(ctx, pinId, userId)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, nil)
}

// Get last 20 users that liked pin
//
//	@Summary		Get last 20 users that liked pin
//	@Description	Get users that liked pin by pin ID
//	@Tags			Likes
//	@Param			pin_id	path	int		true	"Pin ID"
//	@Param			Cookie	header	string	true	"session-token"	default(session-token=)
//	@Produce		json
//	@Success		200	{object}	entity.UserList
//	@Failure		400	{object}	errs.ErrorResponse	"Possible code responses: 12."
//	@Failure		500	{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/likes/{pin_id}/users [get]
func (h *APIHandler) UsersLiked(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	pinId := entity.PinID(id)
	res, errInfo := h.service.GetUsersLiked(ctx, pinId, USERS_LIKED_LIMIT)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteDefaultResponse(w, h.logger, res)
}
