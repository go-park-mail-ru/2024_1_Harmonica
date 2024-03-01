package handler

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"harmonica/db"
	"io"
	"log"
	"net/http"
	"sync"
)

var sessions sync.Map

func (handler *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /login")

	user := new(db.User)

	// Body Collector
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		SetHttpError(w, ErrReadingRequestBody, err, HttpStatusByErr[ErrReadingRequestBody])
		return
	}

	// Body Parser
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		SetHttpError(w, ErrUnmarshalRequestBody, err, 400)
		return
	}

	// Format Validation
	if !ValidateEmail(user.Email) || !ValidatePassword(user.Password) {
		SetHttpError(w, ErrInvalidLoginInput, nil, 400)
		return
	}

}

func (handler *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {}

func (handler *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive POST Request by /register")

	user := new(db.User)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		SetHttpError(w, ErrReadingRequestBody, err, 400)
		return
	}

	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		SetHttpError(w, ErrUnmarshalRequestBody, err, 400)
		return
	}

	if !ValidateEmail(user.Email) || !ValidateNickname(user.Nickname) || !ValidatePassword(user.Password) {
		SetHttpError(w, ErrInvalidRegisterInput, nil, 400)
		return
	}
	// уникальность мэйла и ника проверяется на уровне БД
	// мне кажется тут не надо проверять тогда

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
	if err != nil {
		// позже обработку разных ошибок тут можно будет сделать через switch
		SetHttpError(w, ErrDB, err, 400)
		return
	}
}
