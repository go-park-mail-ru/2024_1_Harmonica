package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Get user by id.
//
//	@Summary		Get user by id
//	@Description	Get user by id in the slug
//	@Tags			Users
//	@Produce		json
//	@Param			Cookie	header		string	true	"session-token"	default(session-token=)
//	@Param			nickname	path		string	true	"User nickname"
//	@Success		200		{object}	entity.UserProfileResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 12, 19."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 7."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			users/{nickname} [get]
func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userNicknameFromSlug := r.PathValue("nickname")
	if !ValidateNickname(userNicknameFromSlug) {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidSlug,
		})
		return
	}

	user, errInfo := h.service.GetUserByNickname(ctx, userNicknameFromSlug)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	if user == emptyUser {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrUserNotExist,
		})
		return
	}

	isOwner := false
	if ctx.Value("is_auth") == true {
		userIdFromSession, ok := ctx.Value("user_id").(entity.UserID)
		if !ok {
			WriteErrorResponse(w, h.logger, errs.ErrorInfo{
				LocalErr: errs.ErrTypeConversion,
			})
		}
		isOwner = user.UserID == userIdFromSession
	}

	userProfile := entity.UserProfileResponse{
		User:           MakeUserResponse(user),
		FollowersCount: 0,
		IsOwner:        isOwner,
	}
	WriteDefaultResponse(w, h.logger, userProfile)

}

// Update user.
//
//	@Summary		Update user
//	@Description	Update user by description and user id.
//	@Tags			Users
//	@Produce		json
//	@Accept			multipart/form-data
//	@Param			Cookie	header		string	true	"session-token"	default(session-token=)
//	@Param			user	formData	entity.User   string	false	"User information in json"
//	@Param			image	formData	file	false	"User avatar"
//	@Success		200		{object}	entity.PinPageResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5, 12, 13, 18"
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 6, 11."
func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIdFromSlug, err := ReadUint64Slug(r, "user_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}

	userIdFromSession, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrTypeConversion,
		})
	}
	if uint64(userIdFromSession) != userIdFromSlug {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrPermissionDenied,
		})
		return
	}

	var user entity.User
	err = UnmarshalRequest(r, &user)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	user.UserID = userIdFromSession // нужно для слоя сервиса! (проверка уникальности ника)

	if user.Nickname != "" && !ValidateNickname(user.Nickname) {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidInputFormat,
		})
		return
	}

	if user.Password != "" {
		if !ValidatePassword(user.Password) {
			WriteErrorResponse(w, h.logger, errs.ErrorInfo{
				LocalErr: errs.ErrInvalidInputFormat,
			})
			return
		}
		hashPassword, errH := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if errH != nil {
			WriteErrorsListResponse(w, h.logger, errs.ErrorInfo{
				GeneralErr: errH,
				LocalErr:   errs.ErrHashingPassword,
			})
			return
		}
		user.Password = string(hashPassword)
	}

	updatedUser, errInfo := h.service.UpdateUser(ctx, user)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}

	WriteDefaultResponse(w, h.logger, MakeUserResponse(updatedUser))
}
