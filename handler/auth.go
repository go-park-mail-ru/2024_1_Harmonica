package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"harmonica/db"
	"harmonica/utils"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var sessions sync.Map

func (handler *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /login")

	isAuth, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, ErrReadCookie)
		return
	}
	if isAuth {
		WriteErrorResponse(w, ErrAlreadyAuthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	userRequest := new(db.User)
	err = json.Unmarshal(bodyBytes, userRequest)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	if !utils.ValidateEmail(userRequest.Email) ||
		!utils.ValidatePassword(userRequest.Password) {
		WriteErrorResponse(w, ErrInvalidInputFormat)
		return
	}

	user, err := handler.connector.GetUserByEmail(userRequest.Email)
	if err != nil {
		WriteErrorResponse(w, ErrDBInternal)
		return
	}
	emptyUser := db.User{}
	if user == emptyUser {
		WriteErrorResponse(w, ErrUserNotExist)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		WriteErrorResponse(w, ErrWrongPassword)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)
	s := session{
		userId: user.UserID,
		expiry: expiresAt,
	}
	sessions.Store(sessionToken, s)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	WriteUserResponse(w, user)

	log.Println("INFO Successful login with session-token:", sessionToken)
}

func (handler *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /logout")

	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			WriteErrorResponse(w, ErrUnauthorized)
			return
		}
		WriteErrorResponse(w, ErrReadCookie)
		return
	}

	sessionToken := c.Value
	s, exists := sessions.Load(sessionToken)
	if !exists || s.(session).IsExpired() {
		WriteErrorResponse(w, ErrUnauthorized)
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
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)

	log.Println("INFO Successful logout")
}

func (handler *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive POST request by /register")

	isAuth, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, ErrReadCookie)
		return
	}
	if isAuth {
		WriteErrorResponse(w, ErrAlreadyAuthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	user := new(db.User)
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		WriteErrorResponse(w, ErrReadingRequestBody)
		return
	}

	if !utils.ValidateEmail(user.Email) ||
		!utils.ValidateNickname(user.Nickname) ||
		!utils.ValidatePassword(user.Password) {
		WriteErrorResponse(w, ErrInvalidInputFormat)
		return
	}
	// уникальность мэйла и ника проверяется на уровне БД

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteErrorResponse(w, ErrHashingPassword)
		return
	}

	foundUser, err := handler.connector.GetUserByEmail(user.Email)
	if err != nil {
		WriteErrorResponse(w, ErrDBInternal)
		return
	}
	emptyUser := db.User{}
	if foundUser != emptyUser {
		WriteErrorResponse(w, ErrUserExist)
		return
	}

	user.Password = string(hashPassword)
	err = handler.connector.RegisterUser(*user)
	var pqErr *pq.Error
	if err != nil && errors.As(err, &pqErr) {
		switch {
		case pqErr.Code == "23505":
			WriteErrorResponse(w, ErrDBUnique)
			return
		default:
			WriteErrorResponse(w, ErrDBInternal)
			return
		}
	}

	// решили, что не авторизуем после реги, а возвращаем 200 просто
	w.WriteHeader(http.StatusOK)

	log.Println("INFO Successful response about authorized user")
}

type session struct {
	userId int64
	expiry time.Time
}

func (handler *APIHandler) IsAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive GET request by /register")
	isAuth, err := CheckAuth(r)
	if err != nil {
		WriteErrorResponse(w, ErrReadCookie)
		return
	}
	if !isAuth {
		WriteErrorResponse(w, ErrUnauthorized)
		return
	}

	c, _ := r.Cookie("session_token")
	// ошибку не проверяю, потому что она уже была бы обработана в CheckAuth
	sessionToken := c.Value
	s, _ := sessions.Load(sessionToken)
	// не проверяю существование ключа, потому что это было обработано в CheckAuth

	user, err := handler.connector.GetUserById(s.(session).userId)
	if err != nil {
		WriteErrorResponse(w, ErrDBInternal)
		return
	}
	emptyUser := db.User{}
	if user == emptyUser {
		WriteErrorResponse(w, ErrUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	WriteUserResponse(w, user)

	log.Println("INFO Authorized user response successfully sent")
}

func (s session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func CheckAuth(r *http.Request) (bool, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return false, nil
		}
		return false, err
	}

	sessionToken := c.Value
	s, exists := sessions.Load(sessionToken)
	if !exists {
		return false, nil
	}
	if s.(session).IsExpired() {
		sessions.Delete(sessionToken)
		return false, nil
	}
	return true, nil
}

func WriteUserResponse(w http.ResponseWriter, user db.User) {
	userResponse := UserResponse{
		UserId:   user.UserID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}
	response, _ := json.Marshal(userResponse)
	w.Write(response)
}
