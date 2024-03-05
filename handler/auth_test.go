package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"

	//"github.com/stretchr/testify/assert"
	"harmonica/db"
	"harmonica/db/mocks"
	"testing"
)

type test struct {
	name               string
	reqBody            map[string]string
	mockBehavior       mockBehavior
	expectedStatusCode int
	expectedRespBody   string
}

type mockBehavior func()

// mock := mocks.NewMethods(t)
var (
	mock         = &mocks.Methods{}
	handler      = APIHandler{Connector: mock}
	loginPath    = "/api/v1/login"
	registerPath = "/api/v1/register"
	logoutPath   = "/api/v1/logout"
	users        = map[string]db.User{
		"200 OK Simple": {
			UserID:   1,
			Email:    "example@gmail.com",
			Nickname: "example",
			Password: "ExampleE1",
		},
		"200 OK Cyrillic alphabet in email": {
			UserID:   2,
			Email:    "москва.тут@москва.рф",
			Nickname: "Moscow",
			Password: "MoscowTest1",
		},
		"200 OK Capital letters in email": {
			UserID:   3,
			Email:    "TesT.her_e@sth.ru",
			Nickname: "Michael",
			Password: "1918Michael",
		},
		"200 OK Crazy (valid) email": {
			UserID:   4,
			Email:    "Crazy%e.{m}ail~@$trange.com",
			Nickname: "crazy_user",
			Password: "crazyUser24",
		},
		"200 OK Max nickname and password length": {
			UserID:   1982897323,
			Email:    "a@b.c",
			Nickname: "crazy_user_2crazy_us",
			Password: "crazyUser24ksHokssn27awb",
		},
	}
)

func PostTypeTestsExecution(t *testing.T, path string, tests []test) {
	for _, curTest := range tests {
		t.Run(curTest.name, func(t *testing.T) {
			reqBytes, err := json.Marshal(curTest.reqBody)
			if err != nil {
				t.Errorf("error marshaling request body: %v", err)
			}
			curTest.mockBehavior()
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", path, bytes.NewBuffer(reqBytes))
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			switch path {
			case loginPath:
				handler.Login(recorder, req)
			case registerPath:
				handler.Register(recorder, req)
			}
			assert.Equal(t, recorder.Code, curTest.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), curTest.expectedRespBody)
		})
	}
}

func TestLoginSuccess(t *testing.T) {
	var tests []test

	for testName, user := range users {
		curHashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			t.Errorf("error hashng passwords: %v", err)
			return
		}
		curTest := test{
			name: testName,
			reqBody: map[string]string{
				"email":    user.Email,
				"password": user.Password,
			},
			mockBehavior: func() {
				user.Password = string(curHashedPassword)
				mock.On("GetUserByEmail", user.Email).Return(user, nil).Once()
			},
			expectedStatusCode: 200,
			expectedRespBody: fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s"}`,
				user.UserID, user.Email, user.Nickname),
		}
		tests = append(tests, curTest)
	}

	PostTypeTestsExecution(t, loginPath, tests)
}

func TestRegisterSuccess(t *testing.T) {
	var tests []test

	for testName, user := range users {
		curTest := test{
			name: testName,
			reqBody: map[string]string{
				"email":    user.Email,
				"nickname": user.Nickname,
				"password": user.Password,
			},
			mockBehavior: func() {
				mock.On("RegisterUser", mock2.AnythingOfType("db.User")).Return(nil)
				mock.On("GetUserByEmail", user.Email).Return(user, nil)
			},
			expectedStatusCode: 200,
			expectedRespBody: fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s"}`,
				user.UserID, user.Email, user.Nickname),
		}
		tests = append(tests, curTest)
	}

	PostTypeTestsExecution(t, registerPath, tests)
}
