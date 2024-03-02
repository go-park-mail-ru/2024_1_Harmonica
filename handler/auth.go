package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"harmonica/db"
	"io"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var sessions sync.Map

type session struct {
	userId int64
	expiry time.Time
}

func (handler *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /login")

	// Checking Existing Authorization
	isAuthorized, err := IsAuthorized(r)
	if err != nil {
		SetHttpError(w, ErrCheckSession, err, HttpStatus[ErrCheckSession])
		return
	}
	if isAuthorized {
		// что делать в случае, если авторизован? послать 403 норм?
		SetHttpError(w, ErrAlreadyAuthorized, err, HttpStatus[ErrAlreadyAuthorized])
		return
	}

	user := new(db.User)

	// Body Collector
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		SetHttpError(w, ErrReadingRequestBody, err, HttpStatus[ErrReadingRequestBody])
		return
	}

	// Body Parser
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		SetHttpError(w, ErrUnmarshalRequestBody, err, HttpStatus[ErrUnmarshalRequestBody])
		return
	}

	// Format Validation
	if !ValidateEmail(user.Email) || !ValidatePassword(user.Password) {
		SetHttpError(w, ErrInvalidLoginInput, nil, HttpStatus[ErrInvalidLoginInput])
		return
	}

	foundUser, err := handler.connector.GetUserByEmail(user.Email)
	// норм просто internal ошибку бд слать тут (500)?
	if err != nil {
		SetHttpError(w, ErrDBInternal, err, HttpStatus[ErrDBInternal])
		return
	}
	if reflect.DeepEqual(foundUser, db.User{}) {
		SetHttpError(w, ErrUserNotExist, nil, HttpStatus[ErrUserNotExist])
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		SetHttpError(w, ErrWrongPassword, err, HttpStatus[ErrWrongPassword])
		return
	}

	// пока поставила просто длительную сессию, но можно сделать короткую + refresh. как лучше?
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	s := session{
		userId: foundUser.UserID,
		expiry: expiresAt,
	}
	sessions.Store(sessionToken, s)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	log.Println("INFO Successful login with session-token:", sessionToken)
}

func (handler *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive POST request by /logout") // POST или GET все-таки?

	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			SetHttpError(w, ErrUnauthorized, err, HttpStatus[ErrUnauthorized])
			return
		}
		SetHttpError(w, nil, err, http.StatusBadRequest)
		return
	}

	sessionToken := c.Value
	s, exists := sessions.Load(sessionToken)
	// сомнения по поводу обработки этих ошибок
	if !exists || s.(session).IsExpired() {
		SetHttpError(w, ErrUnauthorized, err, HttpStatus[ErrUnauthorized])
		return
	}

	//sessions.Delete(sessionToken)
	userId := s.(session).userId
	sessions.Range(func(key, value interface{}) bool {
		userSession := value.(session)
		if userSession.userId == userId {
			sessionToken = key.(string)
			sessions.Delete(sessionToken)
		}
		return true
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	log.Println("INFO Successful logout")
}

func (handler *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive POST request by /register")

	// Checking Existing Authorization
	isAuthorized, err := IsAuthorized(r)
	if err != nil {
		SetHttpError(w, ErrCheckSession, err, HttpStatus[ErrCheckSession])
		return
	}
	if isAuthorized {
		// что делать в случае, если авторизован? послать 403 норм?
		SetHttpError(w, ErrAlreadyAuthorized, err, HttpStatus[ErrAlreadyAuthorized])
		return
	}

	user := new(db.User)

	// Body Collector
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		SetHttpError(w, ErrReadingRequestBody, err, 400)
		return
	}

	// Body Parser
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		SetHttpError(w, ErrUnmarshalRequestBody, err, 400)
		return
	}

	// Format Validation
	if !ValidateEmail(user.Email) || !ValidateNickname(user.Nickname) || !ValidatePassword(user.Password) {
		SetHttpError(w, ErrInvalidRegisterInput, nil, 400)
		return
	}
	// уникальность мэйла и ника проверяется на уровне БД

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		SetHttpError(w, ErrHashingPassword, err, 400)
		return
	}

	foundUser, err := handler.connector.GetUserByEmail(user.Email)
	emptyUser := db.User{}
	if foundUser != emptyUser {
		SetHttpError(w, ErrUserExist, nil, 400)
		return
	}

	user.Password = string(hashPassword)
	err = handler.connector.RegisterUser(*user)
	var pqErr *pq.Error
	if err != nil && errors.As(err, &pqErr) {
		switch {
		case pqErr.Code == "23505":
			SetHttpError(w, ErrDBUnique, err, HttpStatus[ErrDBUnique])
			return
		default:
			SetHttpError(w, ErrDBInternal, err, HttpStatus[ErrDBInternal])
			return
		}
	}

	// после регистрации сразу же авторизуем ?
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	s := session{
		userId: foundUser.UserID,
		expiry: expiresAt,
	}
	sessions.Store(sessionToken, s)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	log.Println("INFO Successful registration and auth with session-token:", sessionToken)
}

func (s session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func IsAuthorized(r *http.Request) (bool, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return false, nil
		}
		return false, err
	}

	sessionToken := c.Value
	userSession, exists := sessions.Load(sessionToken)
	if !exists {
		return false, nil
	}
	if userSession.(session).IsExpired() {
		sessions.Delete(sessionToken)
		return false, nil
	}
	return true, nil
}
