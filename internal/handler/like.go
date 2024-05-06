package handler

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/microservices/like/proto"
	"net/http"

	"google.golang.org/grpc/metadata"
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
	res, err := h.LikeService.SetLike(metadata.NewOutgoingContext(r.Context(),
		metadata.Pairs("request_id", requestId)), &proto.MakeLikeRequest{PinId: int64(pinId), UserId: int64(userId)})
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.ErrGRPCWentWrong})
		return
	}
	if !res.Valid {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.GetLocalErrorByCode[res.LocalError]})
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
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("request_id", requestId))

	pinId, userId, errInfo := GetPinAndUserId(r, ctx)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	res, err := h.LikeService.ClearLike(ctx, &proto.MakeLikeRequest{PinId: int64(pinId), UserId: int64(userId)})
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.ErrGRPCWentWrong})
		return
	}
	if !res.Valid {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.GetLocalErrorByCode[res.LocalError]})
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
	res, err := h.LikeService.GetUsersLiked(metadata.NewOutgoingContext(r.Context(),
		metadata.Pairs("request_id", requestId)), &proto.GetUsersLikedRequest{PinId: int64(pinId), Limit: USERS_LIKED_LIMIT})
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.ErrGRPCWentWrong})
		return
	}
	if !res.Valid {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.GetLocalErrorByCode[res.LocalError]})
		return
	}
	ans := FromGRPCResponseToUsersList(res)
	WriteDefaultResponse(w, h.logger, ans)
}

func FromGRPCResponseToUsersList(protoRes *proto.GetUsersLikedResponse) entity.UserList {
	res := entity.UserList{}
	for _, user := range protoRes.Users {
		res.Users = append(res.Users, entity.UserResponse{
			UserId:    entity.UserID(user.UserId),
			Email:     user.Email,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarUrl,
		})
	}
	return res
}

// Get feed of favorite pins by page.
//
//	@Summary		Get feed of favorite pins by page
//	@Tags			Likes
//	@Param			Cookie	header	string	true	"session-token"	default(session-token=)
//	@Produce		json
//	@Success		200	{object}	entity.FeedPins
//	@Router			/favorites [get]
func (h *APIHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)
	userId := ctx.Value("user_id").(entity.UserID)
	limit, offset, err := GetLimitAndOffset(r)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}
	res, err := h.LikeService.GetFavorites(metadata.NewOutgoingContext(r.Context(),
		metadata.Pairs("request_id", requestId)),
		&proto.GetFavoritesRequest{UserId: int64(userId), Limit: int64(limit), Offset: int64(offset)})

	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.ErrGRPCWentWrong})
		return
	}
	if !res.Valid {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.GetLocalErrorByCode[res.LocalError]})
		return
	}
	feed := FromGRPCResponseToFavoritesFeed(res)
	WriteDefaultResponse(w, h.logger, feed)
}

func FromGRPCResponseToFavoritesFeed(protoRes *proto.GetFavoritesResponse) entity.FeedPins {
	res := entity.FeedPins{}
	for _, pin := range protoRes.Pins {
		res.Pins = append(res.Pins, entity.FeedPinResponse{
			PinId:      entity.PinID(pin.PinId),
			ContentUrl: pin.ContentUrl,
			ContentDX:  pin.ContentDX,
			ContentDY:  pin.ContentDY,
			PinAuthor: entity.PinAuthor{
				UserId:    entity.UserID(pin.Author.UserId),
				Nickname:  pin.Author.Nickname,
				AvatarURL: pin.Author.AvatarUrl,
			},
		})
	}
	return res
}
