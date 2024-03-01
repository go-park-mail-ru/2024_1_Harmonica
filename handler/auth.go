package handler

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"harmonica/db"
	"io"
	"log"
	"net/http"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func (handler *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO receive POST request by /login")

	user := new(db.User)

	// Body Collector
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("here", err)
		http.Error(w, err.Error(), 400)
		return
	}

	// Body Parser
	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		log.Printf("FAIL can't unmarshal data at /login with %s", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	// Validation
	if !ValidateEmail(user.Email) || user.Password == "" {
		log.Printf("FAIL invalid email")
		http.Error(w, "invalid email", 400)
		return
	}
}

func (handler *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {}

func (handler *APIHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO Receive POST Request by /register")

	user := new(db.User)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading the body with %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		log.Printf("here 2 %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !ValidateEmail(user.Email) || !ValidateNickname(user.Nickname) || !ValidatePassword(user.Password) {
		log.Printf("FAIL failed to register (invalid input format)")
		http.Error(w, "failed to register (invalid input format)", http.StatusBadRequest)
		return
	}
	//if !validateNickname(user.Nickname) {
	//	log.Printf("FAIL incorrect nickname format")
	//	http.Error(w, "invalid nickname", http.StatusBadRequest)
	//	return
	//}
	//if {
	//	log.Printf("FAIL password is too long or too short")
	//	http.Error(w, "password is too long or too short", http.StatusBadRequest)
	//	return
	//}
	// уникальность мэйла и ника проверяется на уровне БД
	// мне кажется тут не надо проверять тогда

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	foundUser, err := handler.connector.GetUserByEmail(user.Email)
	emptyUser := db.User{}
	if foundUser != emptyUser {
		log.Printf("FAIL user already exists")
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	user.Password = string(hashPassword)
	err = handler.connector.RegisterUser(*user)
	if err != nil {
		// позже обработку разных ошибок тут можно будет сделать через switch
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
