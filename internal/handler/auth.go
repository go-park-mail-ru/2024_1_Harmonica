package handler

import (
	"encoding/json"
	"errors"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"io"
	"log"
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
func (handler *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /login")
	ctx := r.Context()

	sessionToken, _, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadCookie)
		return
	}
	isAuth := sessionToken != ""
	if isAuth {
		WriteErrorResponse(w, errs.ErrAlreadyAuthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadingRequestBody)
		return
	}

	if !ValidateEmail(user.Email) ||
		!ValidatePassword(user.Password) {
		WriteErrorResponse(w, errs.ErrInvalidInputFormat)
		return
	}

	loggedInUser, err := handler.service.GetUserByEmail(ctx, user.Email)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	if loggedInUser == emptyUser {
		WriteErrorResponse(w, errs.ErrUserNotExist)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(loggedInUser.Password), []byte(user.Password))
	if err != nil {
		WriteErrorResponse(w, errs.ErrWrongPassword)
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
	WriteUserResponse(w, loggedInUser)
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
func (handler *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /logout")

	sessionToken, _, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, errs.ErrReadCookie)
		return
	}
	isAuth := sessionToken != ""
	if !isAuth {
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
//	@Param			request	body		repository.User	true	"json"
//	@Success		200		{object}	entity.UserResponse
//	@Failure		400		{object}	entity.ErrorsListResponse
//	@Failure		403		{object}	entity.ErrorsListResponse
//	@Failure		500		{object}	entity.ErrorsListResponse
//	@Router			/register [post]
func (handler *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive POST request by /users")
	ctx := r.Context()

	sessionToken, _, err := CheckAuth(r)
	if err != nil {
		WriteErrorsListResponse(w, errs.ErrReadCookie)
		return
	}
	isAuth := sessionToken != ""
	if isAuth {
		WriteErrorsListResponse(w, errs.ErrAlreadyAuthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorsListResponse(w, errs.ErrReadingRequestBody)
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorsListResponse(w, errs.ErrReadingRequestBody)
		return
	}

	if !ValidateEmail(user.Email) ||
		!ValidateNickname(user.Nickname) ||
		!ValidatePassword(user.Password) {
		WriteErrorsListResponse(w, errs.ErrInvalidInputFormat)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteErrorsListResponse(w, errs.ErrHashingPassword)
		return
	}
	user.Password = string(hashPassword)
	errsList := handler.service.RegisterUser(ctx, user)
	if len(errsList) != 0 {
		WriteErrorsListResponse(w, errsList...)
		return
	}

	// Search for user by email (to get user id)
	registeredUser, err := handler.service.GetUserByEmail(ctx, user.Email)
	if err != nil {
		WriteErrorsListResponse(w, errs.ErrDBInternal)
		return
	}
	if registeredUser == emptyUser {
		WriteErrorsListResponse(w, errs.ErrUserNotExist)
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
	WriteUserResponse(w, registeredUser)
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
func (handler *APIHandler) IsAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /is_auth")
	ctx := r.Context()

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

	// Checking the existence of user with userId associated with session
	user, err := handler.service.GetUserById(ctx, userIdFromSession)
	if err != nil {
		WriteErrorResponse(w, errs.ErrDBInternal)
		return
	}
	if user == emptyUser {
		WriteErrorResponse(w, errs.ErrUnauthorized)
		return
	}

	WriteUserResponse(w, user)
}

func CheckAuth(r *http.Request) (string, int64, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", 0, nil
		}
		return "", 0, nil
	}
	sessionToken := c.Value
	s, exists := Sessions.Load(sessionToken)
	if !exists {
		return "", 0, nil
	}
	if s.(Session).IsExpired() {
		Sessions.Delete(sessionToken)
		return "", 0, nil
	}
	return sessionToken, s.(Session).UserId, nil
	// в мидлваре прокинуть session_token в контекст, чтобы он был досупен далее в ручке
}

func WriteUserResponse(w http.ResponseWriter, user entity.User) {
	w.Header().Set("Content-Type", "application/json")
	userResponse := entity.UserResponse{
		UserId:   user.UserID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}
	response, _ := json.Marshal(userResponse)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
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
