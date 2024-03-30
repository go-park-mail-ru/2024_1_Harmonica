package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (handler *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	l := handler.logger
	ctx := r.Context()

	userIdFromSlug, err := ReadInt64Slug(r, "user_id")
	if err != nil {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrInvalidSlug,
		})
		return
	}

	userIdFromSession := ctx.Value("user_id").(entity.UserID)
	if uint64(userIdFromSession) != userIdFromSlug {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			LocalErr: errs.ErrDiffUserId,
		})
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	user.UserID = userIdFromSession

	if user.Nickname != "" && !ValidateNickname(user.Nickname) {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidInputFormat,
		})
		return
	}

	if user.Password != "" {
		if !ValidatePassword(user.Password) {
			WriteErrorResponse(w, l, errs.ErrorInfo{
				LocalErr: errs.ErrInvalidInputFormat,
			})
			return
		}
		hashPassword, errH := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if errH != nil {
			WriteErrorsListResponse(w, l, errs.ErrorInfo{
				GeneralErr: errH,
				LocalErr:   errs.ErrHashingPassword,
			})
			return
		}
		user.Password = string(hashPassword)
	}

	updatedUser, errInfo := handler.service.UpdateUser(ctx, user)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, l, errInfo)
		return
	}

	WriteUserResponse(w, l, updatedUser)
}
