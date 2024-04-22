package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func MakeUserResponse(user entity.User) entity.UserResponse {
	userResponse := entity.UserResponse{
		UserId:    user.UserID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}
	return userResponse
}

// Update user.
//
//	@Summary		Update user
//	@Description	Update user by description and user id.
//	@Tags			Users
//	@Produce		json
//	@Accept			json
//	@Param			Cookie		header		string	true	"session-token"	default(session-token=)
//	@Param			nickname	path		string	true	"User nickname"
//	@Success		200			{object}	entity.UserProfileResponse
//	@Failure		400			{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5, 12, 13, 18"
//	@Failure		401			{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403			{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500			{object}	errs.ErrorResponse	"Possible code responses: 6, 11."
//	@Router			/users/{nickname}/ [get]
func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	userNicknameFromSlug := r.PathValue("nickname")
	if !ValidateNickname(userNicknameFromSlug) {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidSlug,
		})
		return
	}
	user, errInfo := h.service.GetUserByNickname(ctx, userNicknameFromSlug)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	if user == emptyUser {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrUserNotExist,
		})
		return
	}

	isOwner := false
	userIdFromSession, ok := ctx.Value("user_id").(entity.UserID)
	if ok {
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
//	@Param			user	formData	string	false	"User information in json"
//	@Param			image	formData	file	false	"User avatar"
//	@Success		200		{object}	entity.PinPageResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5, 12, 13, 18"
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 14."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 6, 11."
//	@Router			/users/{user_id} [post]
func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	userIdFromSlug, err := ReadInt64Slug(r, "user_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrInvalidSlug))
		return
	}

	userIdFromSession, ok := ctx.Value("user_id").(entity.UserID)
	if !ok {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrTypeConversion))
	}
	if uint64(userIdFromSession) != userIdFromSlug {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrDiffUserId))
		return
	}

	var user entity.User

	image, imageHeader, err := r.FormFile("image")
	if err == nil {
		_, name, errUploading := h.service.UploadImage(ctx, image, imageHeader)
		if errUploading != nil {
			WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
				GeneralErr: err,
				LocalErr:   errs.ErrInvalidImg,
			})
			return
		}
		user.AvatarURL = FormImgURL(name)
	}

	userParams := r.FormValue("user")
	err = json.Unmarshal([]byte(userParams), &user)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(err, errs.ErrReadingRequestBody))
		return
	}
	user.UserID = userIdFromSession
	if user.Nickname != "" && !ValidateNickname(user.Nickname) {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
		return
	}

	if user.Password != "" {
		if !ValidatePassword(user.Password) {
			WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrInvalidInputFormat))
			return
		}
		hashPassword, errH := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if errH != nil {
			WriteErrorsListResponse(w, h.logger, requestId, MakeErrorInfo(errH, errs.ErrHashingPassword))
			return
		}
		user.Password = string(hashPassword)
	}

	updatedUser, errInfo := h.service.UpdateUser(ctx, user)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	WriteUserResponse(w, h.logger, updatedUser)
}
