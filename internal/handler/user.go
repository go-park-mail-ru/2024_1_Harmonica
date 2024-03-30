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

	//sessionToken, userIdFromSession, err := CheckAuth(r)
	//if err != nil {
	//	WriteErrorResponse(w, l, errs.ErrorInfo{
	//		GeneralErr: err,
	//		LocalErr:   errs.ErrReadCookie,
	//	})
	//	return
	//}
	//isAuth := sessionToken != ""
	//if !isAuth {
	//	WriteErrorResponse(w, l, errs.ErrorInfo{
	//		//GeneralErr: nil,
	//		LocalErr: errs.ErrUnauthorized,
	//	})
	//	return
	//}

	userIdFromSession := ctx.Value("user_id").(entity.UserID)
	if uint64(userIdFromSession) != userIdFromSlug {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			//GeneralErr: nil,
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

	var user entity.User // раньше создавался указатель тут
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}
	user.UserID = userIdFromSession

	// изменить можно только ник и пароль, остальные поля игнорируются
	if user.Nickname != "" && !ValidateNickname(user.Nickname) {
		WriteErrorResponse(w, l, errs.ErrorInfo{
			//GeneralErr: nil,
			LocalErr: errs.ErrInvalidInputFormat,
		})
		return
	}

	if user.Password != "" {
		if !ValidatePassword(user.Password) {
			WriteErrorResponse(w, l, errs.ErrorInfo{
				//GeneralErr: nil,
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

	WriteUserResponse(w, updatedUser)
}
