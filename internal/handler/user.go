package handler

import (
	"encoding/json"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	auth "harmonica/internal/microservices/auth/proto"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

// Registration
//
//	@Summary		Register user
//	@Description	Register user by POST request and add them to DB
//	@Tags			Authorization
//	@Produce		json
//	@Accept			json
//	@Param			request	body		entity.User	true	"json"
//	@Success		200		{object}	entity.UserResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 7, 8."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 1."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/users [post]
func (h *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorsListResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorsListResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	if !ValidateEmail(user.Email) ||
		!ValidateNickname(user.Nickname) ||
		!ValidatePassword(user.Password) {
		WriteErrorsListResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidInputFormat,
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteErrorsListResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrHashingPassword,
		})
		return
	}
	pass := user.Password
	user.Password = string(hashPassword)

	errsList := h.service.RegisterUser(ctx, user)
	if len(errsList) != 0 {
		WriteErrorsListResponse(w, h.logger, requestId, errsList...)
		return
	}

	// Search for user by email (to get user id)
	registeredUser, errInfo := h.service.GetUserByEmail(ctx, user.Email)
	if errInfo != emptyErrorInfo {
		WriteErrorsListResponse(w, h.logger, requestId, errInfo)
		return
	}
	if registeredUser == emptyUser {
		WriteErrorsListResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrUserNotExist,
		})
		return
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("request_id", requestId))
	res, err := h.AuthService.Login(ctx, &auth.LoginUserRequest{
		UserId:     int64(registeredUser.UserID),
		Email:      registeredUser.Email,
		Nickname:   registeredUser.Nickname,
		AvatarURL:  registeredUser.AvatarURL,
		Password:   pass,
		RegisterAt: registeredUser.RegisterAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.ErrGRPCWentWrong, GeneralErr: err})
		return
	}
	if !res.Valid {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.GetLocalErrorByCode[res.LocalError],
		})
		return
	}
	expiresAt, err := time.Parse(time.RFC3339Nano, res.ExpiresAt)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrCantParseTime,
		})
		return
	}
	SetSessionTokenCookie(w, res.NewSessionToken, expiresAt)
	WriteUserResponse(w, h.logger, registeredUser)
}

// Get user.
//
//	@Summary		Get user
//	@Description	Get user by nickname.
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

	nicknameFromSlug := r.PathValue("nickname")
	var (
		userIdFromSession entity.UserID
		ok                bool
	)
	userId := ctx.Value("user_id")
	if userId != nil {
		userIdFromSession, ok = userId.(entity.UserID)
		if !ok {
			userIdFromSession = 0
		}
	}

	userProfile, errInfo := h.service.GetUserProfileByNickname(ctx, nicknameFromSlug, userIdFromSession)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, requestId, errInfo)
		return
	}
	if (userProfile.User == entity.UserResponse{}) {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrUserNotExist,
		})
		return
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

	_, _, err = r.FormFile("image")
	if err == nil {
		name, errUploading := h.UploadImage(r, "image")
		if errUploading != nil {
			WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
				GeneralErr: err,
				LocalErr:   errs.ErrInvalidImg,
			})
			return
		}
		user.AvatarURL = h.FormImgURL(name)
	}
	userParams := r.FormValue("user")
	err = json.Unmarshal([]byte(userParams), &user)
	fmt.Println(r.FormValue("user"), err)
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

//func MakeUserResponse(user entity.User) entity.UserResponse {
//	userResponse := entity.UserResponse{
//		UserId:    user.UserID,
//		Email:     user.Email,
//		Nickname:  user.Nickname,
//		AvatarURL: user.AvatarURL,
//		AvatarDX:  user.AvatarDX,
//		AvatarDY:  user.AvatarDY,
//	}
//	return userResponse
//}
