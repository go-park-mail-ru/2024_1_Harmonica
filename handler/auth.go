package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"harmonica/db"
	"harmonica/models"
	"harmonica/utils"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	sessions            sync.Map
	sessionTTL          = 24 * time.Hour
	sessionsCleanupTime = 6 * time.Hour
)

func (handler *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /login")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if isAuth {
		WriteErrorResponse(w, ErrAlreadyAuthorized)
		return
	}

	// Body reading
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	// Body parsing
	userRequest := new(db.User)
	err = json.Unmarshal(bodyBytes, userRequest)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	// Format validation
	if !utils.ValidateEmail(userRequest.Email) ||
		!utils.ValidatePassword(userRequest.Password) {
		WriteErrorResponse(w, ErrInvalidInputFormat)
		return
	}

	// Search for user by email
	user, err := handler.connector.GetUserByEmail(userRequest.Email)
	if err != nil {
		WriteErrorResponse(w, ErrDBInternal)
		return
	}
	if user == (db.User{}) {
		WriteErrorResponse(w, ErrUserNotExist)
		return
	}

	// Password check
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		WriteErrorResponse(w, ErrWrongPassword)
		return
	}

	// Session creating
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(sessionTTL)
	s := utils.Session{
		UserId: user.UserID,
		Expiry: expiresAt,
	}
	sessions.Store(sessionToken, s)

	// Writing cookie & response
	SetSessionTokenCookie(w, sessionToken, expiresAt)
	WriteUserResponse(w, user)
}

func (handler *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /logout")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if !isAuth {
		//WriteErrorResponse(w, ErrUnauthorized)
		// решили, что тут обойдемся без ошибки
		w.WriteHeader(http.StatusOK)
		return
	}

	// Deleting the current session
	sessions.Delete(curSessionToken)

	// Writing cookie & response
	SetSessionTokenCookie(w, "", time.Now())
	w.WriteHeader(http.StatusOK)
}

func (handler *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive POST request by /register")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if isAuth {
		WriteErrorResponse(w, ErrAlreadyAuthorized)
		return
	}

	// Body reading
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	// Body parsing
	user := new(db.User)
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	// Format validation
	if !utils.ValidateEmail(user.Email) ||
		!utils.ValidateNickname(user.Nickname) ||
		!utils.ValidatePassword(user.Password) {
		WriteErrorResponse(w, ErrInvalidInputFormat)
		return
	}

	// Password hashing
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteErrorResponse(w, ErrHashingPassword)
		return
	}

	// User registration
	user.Password = string(hashPassword)
	err = handler.connector.RegisterUser(*user)
	var pqErr *pq.Error
	if err != nil && errors.As(err, &pqErr) {
		log.Println(err)
		if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
			WriteErrorResponse(w, ErrDBUniqueEmail)
			return
		}
		if pqErr.Code == "23505" && pqErr.Constraint == "users_nickname_key" {
			WriteErrorResponse(w, ErrDBUniqueNickname)
			return
		}
		WriteErrorResponse(w, ErrDBInternal)
		return
	}

	// Search for user by email (to get user id)
	registeredUser, err := handler.connector.GetUserByEmail(user.Email)
	if err != nil {
		WriteErrorResponse(w, ErrDBInternal)
		return
	}
	if registeredUser == (db.User{}) {
		WriteErrorResponse(w, ErrUserNotExist)
		return
	}

	// Session creating
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(sessionTTL)
	s := utils.Session{
		UserId: registeredUser.UserID,
		Expiry: expiresAt,
	}
	sessions.Store(sessionToken, s)

	// Writing cookie & response
	SetSessionTokenCookie(w, sessionToken, expiresAt)
	WriteUserResponse(w, registeredUser)
}

func (handler *APIHandler) IsAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /is_auth")

	// Checking existing authorization
	curSessionToken, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, ErrReadCookie)
		return
	}
	isAuth := curSessionToken != ""
	if !isAuth {
		WriteErrorResponse(w, ErrUnauthorized)
		return
	}

	// Checking the existence of user with userId associated with session
	s, _ := sessions.Load(curSessionToken)
	// не проверяю существование ключа в мапе, потому что это было обработано в CheckAuth
	user, err := handler.connector.GetUserById(s.(utils.Session).UserId)
	if err != nil {
		WriteErrorResponse(w, ErrDBInternal)
		return
	}
	if user == (db.User{}) {
		WriteErrorResponse(w, ErrUnauthorized)
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
	if s.(utils.Session).IsExpired() {
		sessions.Delete(sessionToken)
		return "", nil
	}
	return sessionToken, nil
}

func WriteUserResponse(w http.ResponseWriter, user db.User) {
	w.Header().Set("Content-Type", "application/json")
	userResponse := models.UserResponse{
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

func CleanupSessions() {
	ticker := time.NewTicker(sessionsCleanupTime)
	for {
		<-ticker.C
		sessions.Range(func(key, value interface{}) bool {
			if session, ok := value.(utils.Session); ok {
				if time.Now().After(session.Expiry) {
					sessions.Delete(key)
				}
			}
			return true
		})
	}
}
