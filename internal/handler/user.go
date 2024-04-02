package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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

	userIdFromSlug, err := ReadInt64Slug(r, "user_id")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}

	userIdFromSession := ctx.Value("user_id").(entity.UserID)
	if uint64(userIdFromSession) != userIdFromSlug {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrDiffUserId,
		})
		return
	}

	var user entity.User

	image, imageHeader, err := r.FormFile("image")
	if err == nil {
		name, errUploading := h.service.UploadImage(ctx, image, imageHeader)
		if errUploading != nil {
			WriteErrorResponse(w, h.logger, errs.ErrorInfo{
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
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	user.UserID = userIdFromSession
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
	WriteUserResponse(w, h.logger, updatedUser)
}
