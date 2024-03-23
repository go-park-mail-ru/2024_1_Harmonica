package handler

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errors_list"
	"io"
	"log"
	"net/http"
)

func (handler *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /update_user")
	ctx := r.Context()

	curSessionToken, curUserId, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if !isAuth {
		WriteErrorResponse(w, errors_list.ErrUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadingRequestBody)
		return
	}

	var user entity.User // раньше создавался указатель тут
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadingRequestBody)
		return
	}
	user.UserID = curUserId

	// изменить можно только ник и пароль, остальное игнорируется
	if user.Nickname != "" && !ValidateNickname(user.Nickname) {
		WriteErrorResponse(w, errors_list.ErrInvalidInputFormat)
		return
	}

	if user.Password != "" {
		if !ValidatePassword(user.Password) {
			WriteErrorResponse(w, errors_list.ErrInvalidInputFormat)
			return
		}
		hashPassword, errH := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if errH != nil {
			WriteErrorsListResponse(w, errors_list.ErrHashingPassword)
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
