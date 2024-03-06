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
	reqBody            interface{}
	mockBehavior       mockBehavior
	expectedStatusCode int
	expectedRespBody   string
}

type mockBehavior func()

var (
	mock         = &mocks.Methods{}
	handler      = APIHandler{Connector: mock}
	loginPath    = "/api/v1/login"
	registerPath = "/api/v1/register"
	logoutPath   = "/api/v1/logout"
	isAuthPath   = "/api/v1/is_auth"
)

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

func TestLoginBad(t *testing.T) {
	var tests []test

	for testName, user := range usersInvalidInput {
		curHashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			t.Errorf("error hashng passwords: %v", err)
			return
		}
		curTest := test{
			name: testName,
			reqBody: map[string]string{
				"email":    user.Email,
				"nickname": user.Nickname,
				"password": user.Password,
			},
			mockBehavior: func() {
				user.Password = string(curHashedPassword)
				mock.On("GetUserByEmail", user.Email).Return(user, nil).Once()
			},
			expectedStatusCode: ErrorCodes[ErrInvalidInputFormat].HttpCode,
			expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
				ErrorCodes[ErrInvalidInputFormat].LocalCode, ErrInvalidInputFormat.Error()),
		}
		tests = append(tests, curTest)
	}

	tests = append(tests, test{
		name:    "Error reading request body",
		reqBody: test{},
		mockBehavior: func() {
			mock.On("GetUserByEmail", "sth").Return(db.User{}, nil)
		},
		expectedStatusCode: ErrorCodes[ErrInvalidInputFormat].HttpCode,
		expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrInvalidInputFormat].LocalCode, ErrInvalidInputFormat.Error()),
	})

	tests = append(tests, test{
		name:    "Error unmarshalling request body",
		reqBody: `"email":"a@b.c,"password":"Test5Name"}`,
		mockBehavior: func() {
			mock.On("GetUserByEmail", "sth").Return(db.User{}, nil)
		},
		expectedStatusCode: ErrorCodes[ErrReadingRequestBody].HttpCode,
		expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrReadingRequestBody].LocalCode, ErrReadingRequestBody.Error()),
	})

	tests = append(tests, test{
		name: "Error find user",
		reqBody: db.User{
			UserID:   3,
			Email:    "Tes.her_e@sth.ru",
			Nickname: "Michael",
			Password: "1918Michael",
		},
		mockBehavior: func() {
			mock.On("GetUserByEmail", "Tes.her_e@sth.ru").Return(db.User{}, nil)
		},
		expectedStatusCode: ErrorCodes[ErrUserNotExist].HttpCode,
		expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrUserNotExist].LocalCode, ErrUserNotExist.Error()),
	})

	wrongHashedPassword, _ := bcrypt.GenerateFromPassword([]byte("WrongPassword"), bcrypt.DefaultCost)
	tests = append(tests, test{
		name: "Error wrong password",
		reqBody: db.User{
			UserID:   3,
			Email:    "Tes.her_e@sth.ru",
			Nickname: "Michael",
			Password: "1918Michael",
		},
		mockBehavior: func() {
			mock.On("GetUserByEmail", "Tes.her_e@sth.ru").Return(db.User{
				UserID:   3,
				Email:    "Tes.her_e@sth.ru",
				Nickname: "Michael",
				Password: string(wrongHashedPassword),
			}, nil)
		},
		expectedStatusCode: ErrorCodes[ErrUserNotExist].HttpCode,
		expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrUserNotExist].LocalCode, ErrUserNotExist.Error()),
	})

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

func TestRegisterBad(t *testing.T) {
	var tests []test

	for testName, user := range usersInvalidInput {
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
			expectedStatusCode: ErrorCodes[ErrInvalidInputFormat].HttpCode,
			expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
				ErrorCodes[ErrInvalidInputFormat].LocalCode, ErrInvalidInputFormat.Error()),
		}
		tests = append(tests, curTest)
	}

	tests = append(tests, test{
		name: "Invalid characters in nickname",
		reqBody: db.User{
			UserID:   3,
			Email:    "Tes.her_e@sth.ru",
			Nickname: "Micha'[el#",
			Password: "1918Michael",
		},
		mockBehavior: func() {
			mock.On("RegisterUser", mock2.AnythingOfType("db.User")).Return(nil)
			mock.On("GetUserByEmail", "sth").Return(db.User{}, nil)
		},
		expectedStatusCode: ErrorCodes[ErrInvalidInputFormat].HttpCode,
		expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrInvalidInputFormat].LocalCode, ErrInvalidInputFormat.Error()),
	})

	tests = append(tests, test{
		name:    "Error reading request body",
		reqBody: true,
		mockBehavior: func() {
			mock.On("RegisterUser", mock2.AnythingOfType("db.User")).Return(nil)
			mock.On("GetUserByEmail", "sth").Return(db.User{}, nil)
		},
		expectedStatusCode: ErrorCodes[ErrReadingRequestBody].HttpCode,
		expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrReadingRequestBody].LocalCode, ErrReadingRequestBody.Error()),
	})

	tests = append(tests, test{
		name:    "Error unmarshalling request body",
		reqBody: `{"email":"a@b.c,"nickname":"sss","password":"Test5Name"}`,
		mockBehavior: func() {
			mock.On("RegisterUser", mock2.AnythingOfType("db.User")).Return(nil)
			mock.On("GetUserByEmail", "sth").Return(db.User{}, nil)
		},
		expectedStatusCode: ErrorCodes[ErrReadingRequestBody].HttpCode,
		expectedRespBody: fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			ErrorCodes[ErrReadingRequestBody].LocalCode, ErrReadingRequestBody.Error()),
	})

	PostTypeTestsExecution(t, registerPath, tests)
}

func TestLogout(t *testing.T) {
	var tests []test

	for testName, _ := range users {
		curTest := test{
			name:               testName,
			reqBody:            map[string]string{},
			mockBehavior:       func() {},
			expectedStatusCode: 200,
			expectedRespBody:   "",
		}
		tests = append(tests, curTest)
	}

	GetTypeTestsExecution(t, logoutPath, tests)
}

func TestIsAuth(t *testing.T) {
	var tests []test

	for testName, _ := range users {
		curTest := test{
			name:               testName,
			reqBody:            map[string]string{},
			mockBehavior:       func() {},
			expectedStatusCode: 401,
			expectedRespBody:   `{"code":2,"message":"unauthorized"}`,
		}
		tests = append(tests, curTest)
	}

	GetTypeTestsExecution(t, isAuthPath, tests)
}

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

func GetTypeTestsExecution(t *testing.T, path string, tests []test) {
	for _, curTest := range tests {
		t.Run(curTest.name, func(t *testing.T) {
			curTest.mockBehavior()
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			switch path {
			case logoutPath:
				handler.Logout(recorder, req)
			case isAuthPath:
				handler.IsAuth(recorder, req)
			}
			assert.Equal(t, recorder.Code, curTest.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), curTest.expectedRespBody)
		})
	}
}

var (
	users = map[string]db.User{
		"200 Capital letters in email": {
			UserID:   3,
			Email:    "TesT.her_e@sth.ru",
			Nickname: "Michael",
			Password: "1918Michael",
		},
		"200 Crazy (valid) email": {
			UserID:   4,
			Email:    "Crazy%e.{m}ail~@$trange.com",
			Nickname: "crazy_user",
			Password: "crazyUser24",
		},
		"200 Max nickname and password length": {
			UserID:   1982897323,
			Email:    "a@b.c",
			Nickname: "crazy_user_2crazy_us",
			Password: "crazyUser24ksHokssn27awb",
		},
	}
	usersInvalidInput = map[string]db.User{
		"400 No @ in email": {
			UserID:   1,
			Email:    "badexamplegmail.com",
			Nickname: "example",
			Password: "ExampleE1",
		},
		"400 Two dots in email": {
			UserID:   2,
			Email:    "москва..тут@москва.рф",
			Nickname: "Moscow",
			Password: "MoscowTest1",
		},
		"400 Invalid Password": {
			UserID:   4,
			Email:    "москва.тут@москва.рф",
			Nickname: "crazy_user",
			Password: "crazyuser24",
		},
	}
)
