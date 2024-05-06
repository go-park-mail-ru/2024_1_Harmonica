package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"harmonica/internal/handler"
	"harmonica/internal/microservices/auth/proto"
	grpcAuth "harmonica/mocks/microservices/auth/proto"
	mock_proto "harmonica/mocks/microservices/auth/proto"
	mock_service "harmonica/mocks/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler_Login(t *testing.T) {
	curTime := time.Time{}
	timeString := curTime.Format(time.RFC3339Nano)
	userJsonString := fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s","avatar_url":"","avatar_width":0,"avatar_height":0}`,
		0, "email@mail.ru", "MIRACLE")
	errorJsonString := fmt.Sprintf(`{"code":11,"message":"internal db error"}`)

	type AuthServiceMockBehavior func(
		mockClient *grpcAuth.MockAuthorizationClient,
		req *proto.LoginUserRequest,
		expectedResponse *proto.LoginUserResponse,
	)

	testTable := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
		expectedJSON       string

		expectedRequest         *proto.LoginUserRequest
		expectedResponse        *proto.LoginUserResponse
		authServiceMockBehavior AuthServiceMockBehavior

		context context.Context
	}{
		{
			name:               "OK test case 1",
			requestBody:        map[string]string{"email": "email@mail.ru", "password": "Passw0rd", "nickname": "MIRACLE"},
			expectedStatusCode: 200,
			expectedJSON:       userJsonString,
			expectedRequest: &proto.LoginUserRequest{
				UserId:   0,
				Email:    "email@mail.ru",
				Nickname: "MIRACLE",
				Password: "Passw0rd",
			},
			expectedResponse: &proto.LoginUserResponse{
				UserId:     0,
				Email:      "email@mail.ru",
				Nickname:   "MIRACLE",
				Password:   "Passw0rd",
				RegisterAt: timeString,
				LocalError: 0,
				Valid:      true,
				AvatarURL:  "",
				ExpiresAt:  timeString,
			},
			authServiceMockBehavior: func(mockClient *mock_proto.MockAuthorizationClient,
				req *proto.LoginUserRequest, expectedResponse *proto.LoginUserResponse) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context: context.WithValue(context.Background(), "request_id", "1"),
		},
		{
			name:               "Error test case 1",
			requestBody:        map[string]string{"email": "email@mail.ru", "password": "Passw0rd", "nickname": "MIRACLE"},
			expectedStatusCode: 500,
			expectedJSON:       errorJsonString,
			expectedRequest: &proto.LoginUserRequest{
				UserId:   0,
				Email:    "email@mail.ru",
				Nickname: "MIRACLE",
				Password: "Passw0rd",
			},
			expectedResponse: &proto.LoginUserResponse{
				UserId:     0,
				Email:      "email@mail.ru",
				Nickname:   "MIRACLE",
				Password:   "Passw0rd",
				RegisterAt: timeString,
				LocalError: 11,
				Valid:      false,
				AvatarURL:  "",
				ExpiresAt:  timeString,
			},
			authServiceMockBehavior: func(mockClient *mock_proto.MockAuthorizationClient,
				req *proto.LoginUserRequest, expectedResponse *proto.LoginUserResponse) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context: context.WithValue(context.Background(), "request_id", "1"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthServiceClient := mock_proto.NewMockAuthorizationClient(ctrl)
			testCase.authServiceMockBehavior(mockAuthServiceClient, testCase.expectedRequest, testCase.expectedResponse)
			serviceMock := mock_service.NewMockIService(ctrl)

			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, mockAuthServiceClient, nil, nil)
			requestJSON, _ := json.Marshal(&testCase.requestBody)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(requestJSON))
			r = r.WithContext(testCase.context)

			h.Login(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}
