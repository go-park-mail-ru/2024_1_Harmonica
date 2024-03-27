package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
)

func GetPinAndUserId(r *http.Request) (entity.PinID, entity.UserID, error) {
	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		return entity.PinID(0), entity.UserID(0), errs.ErrInvalidSlug
	}
	pinId := entity.PinID(id)
	_, userId, err := CheckAuth(r)
	if err != nil || userId == 0 {
		return entity.PinID(0), entity.UserID(0), errs.ErrReadCookie
	}
	return pinId, userId, nil
}

func (handler *APIHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	pinId, userId, err := GetPinAndUserId(r)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	err = handler.service.CreateLike(r.Context(), pinId, userId)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	WriteDefaultResponse(w, nil)
}

func (handler *APIHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	pinId, userId, err := GetPinAndUserId(r)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	err = handler.service.DeleteLike(r.Context(), pinId, userId)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
	WriteDefaultResponse(w, nil)
}

func (handler *APIHandler) LikedPins(w http.ResponseWriter, r *http.Request) {
	_, userId, err := CheckAuth(r)
	if err != nil || userId == 0 {
		WriteErrorResponse(w, errs.ErrReadCookie)
		return
	}
	res, err := handler.service.GetLikedPins(r.Context(), userId, Limit)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	WriteDefaultResponse(w, res)
}

func (handler *APIHandler) LikesCount(w http.ResponseWriter, r *http.Request) {
	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, errs.ErrInvalidSlug)
		return
	}
	pinId := entity.PinID(id)
	count, err := handler.service.GetLikesCount(r.Context(), pinId)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	WriteDefaultResponse(w, map[string]uint64{"likes_count": count})
}

func (handler *APIHandler) UsersLiked(w http.ResponseWriter, r *http.Request) {
	id, err := ReadInt64Slug(r, "pin_id")
	if err != nil {
		WriteErrorResponse(w, errs.ErrInvalidSlug)
		return
	}
	pinId := entity.PinID(id)
	res, err := handler.service.GetUsersLiked(r.Context(), pinId, Limit)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	WriteDefaultResponse(w, res)
}
