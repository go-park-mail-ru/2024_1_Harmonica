package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"

	//"github.com/stretchr/testify/assert"
	"harmonica/db"
	"harmonica/db/mocks"
	"testing"
)

func TestLogin(t *testing.T) {
	mock := mocks.NewMethods(t)
	handler := APIHandler{Connector: mock}
	type mockBehavior func()

	users := []db.User{
		{
			UserID:   123,
			Email:    "example@gmail.com",
			Nickname: "example",
			Password: "ExampleE1",
		},
		{
			UserID:   2,
			Email:    "TesT.here@что-то.рф",
			Nickname: "test1",
			Password: "Testtest1T",
		},
	}

	var hashPasswords []string
	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			t.Errorf("error hashng passwords: %v", err)
			return
		}
		hashPasswords = append(hashPasswords, string(hashedPassword))
	}

	tests := []struct {
		name               string
		reqBody            map[string]string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedRespBody   string
	}{
		{
			name: "OK",
			reqBody: map[string]string{
				"email":    users[0].Email,
				"password": users[0].Password,
			},
			mockBehavior: func() {
				users[0].Password = hashPasswords[0]
				mock.On("GetUserByEmail", users[0].Email).Return(users[0], nil).Once()
			},
			expectedStatusCode: 200,
			expectedRespBody: fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s"}`,
				users[0].UserID, users[0].Email, users[0].Nickname),
		},
		{
			name: "OK",
			reqBody: map[string]string{
				"email":    users[1].Email,
				"password": users[1].Password,
			},
			mockBehavior: func() {
				users[1].Password = hashPasswords[1]
				mock.On("GetUserByEmail", users[1].Email).Return(users[1], nil).Once()
			},
			expectedStatusCode: 200,
			expectedRespBody: fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s"}`,
				users[1].UserID, users[1].Email, users[1].Nickname),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBytes, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Errorf("error marshaling request body: %v", err)
			}
			test.mockBehavior()
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBytes))
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			handler.Login(recorder, req)
			assert.Equal(t, recorder.Code, test.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), test.expectedRespBody)
		})
	}
}
