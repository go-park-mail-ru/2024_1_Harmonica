package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
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
			LocalErr: errs.ErrDiffUserId,
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
