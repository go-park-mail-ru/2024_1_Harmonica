package handler

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	Sessions            sync.Map
	sessionTTL          = 24 * time.Hour
	sessionsCleanupTime = 6 * time.Hour
	emptyUser           = entity.User{}
	emptyErrorInfo      = errs.ErrorInfo{}
)

// Login
//
//	@Summary		Login user
//	@Description	Login user by request.body json
//	@Tags			Authorization
//
// @Param 		 Cookie header string  false "session-token"     default(session-token=)
// @Success		200		{object}	interface{}
// @Failure		400		{object}	errs.ErrorResponse
// @Failure		401		{object}	errs.ErrorResponse
// @Failure		403		{object}	errs.ErrorResponse
// @Failure		500		{object}	errs.ErrorResponse
// @Header			200		{string}	Set-Cookie	"session-token"
// @Router			/login [post]
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User
	err := UnmarshalRequest(r, &user)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody})
		return
	}

	if !ValidateEmail(user.Email) ||
		!ValidatePassword(user.Password) {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidInputFormat,
		})
		return
	}

	loggedInUser, errInfo := h.service.GetUserByEmail(ctx, user.Email)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	if loggedInUser == emptyUser {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrUserNotExist,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(loggedInUser.Password), []byte(user.Password))
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrWrongPassword})
		return
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(sessionTTL)
	s := Session{
		UserId: loggedInUser.UserID,
		Expiry: expiresAt,
	}
	Sessions.Store(newSessionToken, s)

	SetSessionTokenCookie(w, newSessionToken, expiresAt)
	WriteDefaultResponse(w, h.logger, MakeUserResponse(loggedInUser))
}

// Logout
//
//	@Summary		Logout user
//	@Description	Logout user by their session cookie
//	@Tags			Authorization
//
// @Param 		 Cookie header string  true "session-token"     default(session-token=)
//
//	@Success		200		{object}	interface{}
//	@Failure		400		{object}	errs.ErrorResponse
//	@Header			200		{string}	Set-Cookie	"session-token"
//	@Router			/logout [get]
func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ctx.Value("is_auth") == false {
		w.WriteHeader(http.StatusOK)
		return
	}

	sessionToken := ctx.Value("session_token")
	Sessions.Delete(sessionToken)
	SetSessionTokenCookie(w, "", time.Now())
	w.WriteHeader(http.StatusOK)
}

// Registration
//
//	@Summary		Register user
//	@Description	Register user by POST request and add them to DB
//	@Tags			Authorization
//	@Produce		json
//	@Accept			json
//	@Param			request	body		repository.User	true	"json"
//	@Success		200		{object}	entity.UserResponse
//	@Failure		400		{object}	entity.ErrorsListResponse
//	@Failure		403		{object}	entity.ErrorsListResponse
//	@Failure		500		{object}	entity.ErrorsListResponse
//	@Router			/register [post]
func (h *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user entity.User
	err := UnmarshalRequest(r, &user)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody})
		return
	}

	if !ValidateEmail(user.Email) ||
		!ValidateNickname(user.Nickname) ||
		!ValidatePassword(user.Password) {
		WriteErrorsListResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrInvalidInputFormat,
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteErrorsListResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrHashingPassword,
		})
		return
	}
	user.Password = string(hashPassword)

	errsList := h.service.RegisterUser(ctx, user)
	if len(errsList) != 0 {
		WriteErrorsListResponse(w, h.logger, errsList...)
		return
	}

	// Search for user by email (to get user id)
	registeredUser, errInfo := h.service.GetUserByEmail(ctx, user.Email)
	if errInfo != emptyErrorInfo {
		WriteErrorsListResponse(w, h.logger, errInfo)
		return
	}
	if registeredUser == emptyUser {
		WriteErrorsListResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrUserNotExist,
		})
		return
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(sessionTTL)
	s := Session{
		UserId: registeredUser.UserID,
		Expiry: expiresAt,
	}
	Sessions.Store(newSessionToken, s)

	SetSessionTokenCookie(w, newSessionToken, expiresAt)
	WriteDefaultResponse(w, h.logger, MakeUserResponse(registeredUser))
}

// Check if user is authorized
//
//	@Summary		Get auth status
//	@Description	Get user by request cookie
//	@Tags			Authorization
//
// @Param 		 Cookie header string  false "session-token"     default(session-token=)
//
//	@Produce		json
//	@Success		200	{object}	entity.UserResponse
//	@Failure		400	{object}	errs.ErrorResponse
//	@Failure		401	{object}	errs.ErrorResponse
//	@Failure		500	{object}	errs.ErrorResponse
//	@Router			/is_auth [get]
func (h *APIHandler) IsAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIdFromSession := ctx.Value("user_id").(entity.UserID)

	// Checking the existence of user with userId associated with session
	user, errInfo := h.service.GetUserById(ctx, userIdFromSession)
	if errInfo != emptyErrorInfo {
		WriteErrorResponse(w, h.logger, errInfo)
		return
	}
	if user == emptyUser {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			LocalErr: errs.ErrUnauthorized})
		return
	}

	WriteDefaultResponse(w, h.logger, MakeUserResponse(user))
}
