package handler

import (
	"encoding/json"
	"errors"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

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
//	@Produce		json
//	@Accept			json
//	@Header			200		{string}	Set-Cookie	"session-token"
//	@Param			Cookie	header		string		false	"session-token"	default(session-token=)
//	@Success		200		{object}	entity.User
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 7, 8."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 1."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/login [post]
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
	WriteUserResponse(w, h.logger, loggedInUser)
}

// Logout
//
//	@Summary		Logout user
//	@Description	Logout user by their session cookie
//	@Tags			Authorization
//	@Param			Cookie	header		string		false	"session-token"	default(session-token=)
//	@Header			200		{string}	Set-Cookie	"session-token"
//	@Success		200
//	@Failure		400	{object}	errs.ErrorResponse	"Possible code responses: 3."
//	@Router			/logout [get]
func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusOK)
			return
		}
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadCookie,
		})
		return
	}
	sessionToken := c.Value
	_, exists := Sessions.Load(sessionToken)
	if !exists {
		SetSessionTokenCookie(w, "", time.Now())
		w.WriteHeader(http.StatusOK)
		return
	}

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
//	@Param			request	body		entity.User	true	"json"
//	@Success		200		{object}	entity.UserResponse
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 7, 8."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 1."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/users [post]
func (h *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorsListResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorsListResponse(w, h.logger, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
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
	WriteUserResponse(w, h.logger, registeredUser)
}

// Check if user is authorized
//
//	@Summary		Get auth status
//	@Description	Get user by request cookie
//	@Tags			Authorization
//	@Param			Cookie	header	string	false	"session-token"	default(session-token=)
//	@Produce		json
//	@Success		200	{object}	entity.UserResponse
//	@Failure		400	{object}	errs.ErrorResponse	"Possible code responses: 3."
//	@Failure		401	{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		500	{object}	errs.ErrorResponse	"Possible code responses: 11."
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

	WriteUserResponse(w, h.logger, user)
}

func WriteUserResponse(w http.ResponseWriter, logger *zap.Logger, user entity.User) {
	w.Header().Set("Content-Type", "application/json")
	userResponse := entity.UserResponse{
		UserId:    user.UserID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}
	response, _ := json.Marshal(userResponse)
	_, err := w.Write(response)
	if err != nil {
		logger.Error(
			errs.ErrServerInternal.Error(),
			zap.Int("local_error_code", errs.ErrorCodes[errs.ErrServerInternal].LocalCode),
			zap.String("general_error", err.Error()),
		)
	}
}

func SetSessionTokenCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
	})
}
