package handler

import (
	"encoding/json"
	"errors"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errors_list"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	sessions            sync.Map
	sessionTTL          = 24 * time.Hour
	sessionsCleanupTime = 6 * time.Hour
)

// Login
//
//	@Summary		Login user
//	@Description	Login user by request.body json
//	@Tags			Authorization
//
// @Param 		 Cookie header string  false "session-token"     default(session-token=)
// @Success		200		{object}	interface{}
// @Failure		400		{object}	entity.ErrorResponse
// @Failure		401		{object}	entity.ErrorResponse
// @Failure		403		{object}	entity.ErrorResponse
// @Failure		500		{object}	entity.ErrorResponse
// @Header			200		{string}	Set-Cookie	"session-token"
// @Router			/login [post]
func (handler *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context() (чтобы session-token можно было получить после мидлвары)
	log.Println("INFO receive POST request by /login")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if isAuth {
		WriteErrorResponse(w, errors_list.ErrAlreadyAuthorized)
		return
	}

	// Body reading
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadingRequestBody)
		return
	}

	// Body parsing
	userRequest := new(entity.User)
	err = json.Unmarshal(bodyBytes, userRequest)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadingRequestBody)
		return
	}

	// Format validation
	if !ValidateEmail(userRequest.Email) ||
		!ValidatePassword(userRequest.Password) {
		WriteErrorResponse(w, errors_list.ErrInvalidInputFormat)
		return
	}

	// Search for user by email
	user, err := handler.service.GetUserByEmail(userRequest.Email)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrDBInternal)
		return
	}
	if user == (entity.User{}) {
		WriteErrorResponse(w, errors_list.ErrUserNotExist)
		return
	}

	// Password check
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrWrongPassword)
		return
	}

	// Session creating
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(sessionTTL)
	s := Session{
		UserId: user.UserID,
		Expiry: expiresAt,
	}
	sessions.Store(sessionToken, s)

	// Writing cookie & response
	SetSessionTokenCookie(w, sessionToken, expiresAt)
	WriteUserResponse(w, user)
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
//	@Failure		400		{object}	entity.ErrorResponse
//	@Header			200		{string}	Set-Cookie	"session-token"
//	@Router			/logout [get]
func (handler *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /logout")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if !isAuth {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Deleting the current session
	sessions.Delete(curSessionToken)

	// Writing cookie & response
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
	log.Println("INFO Receive POST request by /register")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorsListResponse(w, errors_list.ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if isAuth {
		WriteErrorsListResponse(w, errors_list.ErrAlreadyAuthorized)
		return
	}

	// Body reading
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorsListResponse(w, errors_list.ErrReadingRequestBody)
		return
	}

	// Body parsing
	user := new(entity.User)
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		WriteErrorsListResponse(w, errors_list.ErrReadingRequestBody)
		return
	}

	// Format validation
	if !ValidateEmail(user.Email) ||
		!ValidateNickname(user.Nickname) ||
		!ValidatePassword(user.Password) {
		WriteErrorsListResponse(w, errors_list.ErrInvalidInputFormat)
		return
	}

	// Password hashing
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteErrorsListResponse(w, errors_list.ErrHashingPassword)
		return
	}

	// User registration (checking for unique fields is now here)
	user.Password = string(hashPassword)
	errs := handler.service.RegisterUser(*user)
	if len(errs) != 0 {
		WriteErrorsListResponse(w, errs...)
		return
	}

	// Search for user by email (to get user id)
	registeredUser, err := handler.service.GetUserByEmail(user.Email)
	if err != nil {
		WriteErrorsListResponse(w, errors_list.ErrDBInternal)
		return
	}
	if registeredUser == (entity.User{}) {
		WriteErrorsListResponse(w, errors_list.ErrUserNotExist)
		return
	}

	// Session creating
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(sessionTTL)
	s := Session{
		UserId: registeredUser.UserID,
		Expiry: expiresAt,
	}
	sessions.Store(sessionToken, s)

	// Writing cookie & response
	SetSessionTokenCookie(w, sessionToken, expiresAt)
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
//	@Failure		400	{object}	entity.ErrorResponse
//	@Failure		401	{object}	entity.ErrorResponse
//	@Failure		500	{object}	entity.ErrorResponse
//	@Router			/is_auth [get]
func (handler *APIHandler) IsAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /is_auth")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if !isAuth {
		WriteErrorResponse(w, errors_list.ErrUnauthorized)
		return
	}

	// Checking the existence of user with userId associated with session
	s, _ := sessions.Load(curSessionToken)
	// не проверяю существование ключа в мапе, потому что это было обработано в CheckAuth
	user, err := handler.service.GetUserById(s.(Session).UserId)
	if err != nil {
		WriteErrorResponse(w, errors_list.ErrDBInternal)
		return
	}
	if user == (entity.User{}) {
		WriteErrorResponse(w, errors_list.ErrUnauthorized)
		return
	}

	// Writing response
	WriteUserResponse(w, user)
}

func CheckAuth(r *http.Request) (string, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		return "", err
	}
	sessionToken := c.Value
	s, exists := sessions.Load(sessionToken)
	if !exists {
		return "", nil
	}
	if s.(Session).IsExpired() {
		sessions.Delete(sessionToken)
		return "", nil
	}
	return sessionToken, nil
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
