package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (handler *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /users/{user_id}")
	ctx := r.Context()

	userIdFromSlug, err := ReadInt64Slug(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrInvalidSlug)
		return
	}

	sessionToken, userIdFromSession, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadCookie)
		return
	}
	isAuth := sessionToken != ""
	if !isAuth {
		WriteErrorResponse(w, errs.ErrUnauthorized)
		return
	}

	if uint64(userIdFromSession) != userIdFromSlug {
		WriteErrorResponse(w, errs.ErrDiffUserId)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}

	var user entity.User // раньше создавался указатель тут
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}
	user.UserID = userIdFromSession

	// изменить можно только ник и пароль, остальные поля игнорируются
	if user.Nickname != "" && !ValidateNickname(user.Nickname) {
		WriteErrorResponse(w, errs.ErrInvalidInputFormat)
		return
	}

	if user.Password != "" {
		if !ValidatePassword(user.Password) {
			WriteErrorResponse(w, errs.ErrInvalidInputFormat)
			return
		}
		hashPassword, errH := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if errH != nil {
			WriteErrorsListResponse(w, errs.ErrHashingPassword)
			return
		}
		user.Password = string(hashPassword)
	}

	updatedUser, err := handler.service.UpdateUser(ctx, user)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}

	WriteUserResponse(w, updatedUser)
}
