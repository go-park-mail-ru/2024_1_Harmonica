package handler

import (
	"encoding/json"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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
		userIdFromSession := ctx.Value("user_id").(entity.UserID)
		isOwner = user.UserID == userIdFromSession
	}

	userProfile := entity.UserProfileResponse{
		User:            MakeUserResponse(user),
		FollowersNumber: 0,
		IsOwner:         isOwner,
	}
	WriteDefaultResponse(w, h.logger, userProfile)

}

func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIdFromSlug, err := ReadUint64Slug(r, "user_id")
	fmt.Println(userIdFromSlug)
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

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)
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
	//WriteUserResponse(w, h.logger, updatedUser)
}
