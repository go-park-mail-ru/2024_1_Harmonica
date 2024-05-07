package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
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

var users = []entity.User{
	{
		UserID:   1,
		Email:    "mary@email.ru",
		Nickname: "MaryPoppins",
		Password: "MaryPoppins25",
	},
	{
		UserID:   2,
		Email:    "TesT.her_e@sth.ru",
		Nickname: "Michael",
		Password: "1918Michael",
	},
	{
		UserID:   3,
		Email:    "Crazy%e.{m}ail~@$trange.com",
		Nickname: "crazy_user",
		Password: "crazyUser24ksHokssn27awb",
	},
	{
		UserID:   4,
		Email:    "something_wrong",
		Nickname: "crazy_user",
		Password: "password",
	},
}

func MakeUserResponseBody(user entity.User) string {
	return fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s","avatar_url":"%s"}`,
		user.UserID, user.Email, user.Nickname, user.AvatarURL)
}
func MakeUserResponseBodyWithBounds(user entity.User) string {
	return fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s","avatar_url":"%s","avatar_width":0,"avatar_height":0}`,
		user.UserID, user.Email, user.Nickname, user.AvatarURL)
}

func MakeErrorResponse(err error) string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`,
		errs.ErrorCodes[err].LocalCode, err.Error())
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func MakeErrorListResponse(errsList ...error) string {
	errors := make([]Error, len(errsList))
	for i, err := range errsList {
		errors[i] = Error{
			Code:    errs.ErrorCodes[err].LocalCode,
			Message: err.Error(),
		}
	}
	response := ErrorResponse{Errors: errors}
	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

func TestHandler_Login(t *testing.T) {
	curTime := time.Time{}
	timeString := curTime.Format(time.RFC3339Nano)
	userJsonString := fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s","avatar_url":"","avatar_width":0,"avatar_height":0}`,
		0, "email@mail.ru", "MIRACLE")

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
			expectedJSON:       MakeErrorResponse(errs.ErrDBInternal),
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
		{
			name:               "Error test case 2",
			requestBody:        map[string]string{"email": "email@mail.ru", "password": "Passw0rd", "nickname": "MIRACLE"},
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrCantParseTime),
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
				ExpiresAt:  "9 часов утра, 07.05.24",
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

func TestHandler_IsAuth(t *testing.T) {
	userJsonString := fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s","avatar_url":"","avatar_width":0,"avatar_height":0}`,
		0, "email@mail.ru", "MIRACLE")
	errorJsonStringDBErr := `{"code":11,"message":"internal db error"}`

	type AuthServiceMockBehavior func(
		mockClient *grpcAuth.MockAuthorizationClient,
		req *proto.Empty,
		expectedResponse *proto.IsAuthResponse,
	)

	testTable := []struct {
		name               string
		expectedStatusCode int
		expectedJSON       string

		expectedRequest         *proto.Empty
		expectedResponse        *proto.IsAuthResponse
		authServiceMockBehavior AuthServiceMockBehavior

		context context.Context
	}{
		{
			name:               "OK test case 1",
			expectedStatusCode: 200,
			expectedJSON:       userJsonString,
			expectedRequest:    &proto.Empty{},
			expectedResponse: &proto.IsAuthResponse{
				IsAuthorized: true,
				User: &proto.IsAuthUserResponse{
					UserId:    0,
					Email:     "email@mail.ru",
					Nickname:  "MIRACLE",
					AvatarURL: "",
				},
				LocalError: 0,
				Valid:      true,
			},
			authServiceMockBehavior: func(mockClient *mock_proto.MockAuthorizationClient,
				req *proto.Empty, expectedResponse *proto.IsAuthResponse) {
				mockClient.EXPECT().IsAuth(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "1"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test case 1",
			expectedStatusCode: 500,
			expectedJSON:       errorJsonStringDBErr,
			expectedRequest:    &proto.Empty{},
			expectedResponse: &proto.IsAuthResponse{
				IsAuthorized: true,
				User: &proto.IsAuthUserResponse{
					UserId:    0,
					Email:     "email@mail.ru",
					Nickname:  "MIRACLE",
					AvatarURL: "",
				},
				LocalError: 11,
				Valid:      false,
			},
			authServiceMockBehavior: func(mockClient *mock_proto.MockAuthorizationClient,
				req *proto.Empty, expectedResponse *proto.IsAuthResponse) {
				mockClient.EXPECT().IsAuth(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "1"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test case 2",
			expectedStatusCode: 500,
			expectedJSON:       errorJsonStringDBErr,
			expectedRequest:    &proto.Empty{},
			expectedResponse: &proto.IsAuthResponse{
				IsAuthorized: true,
				User: &proto.IsAuthUserResponse{
					UserId:    0,
					Email:     "email@mail.ru",
					Nickname:  "MIRACLE",
					AvatarURL: "",
				},
				LocalError: 11,
				Valid:      false,
			},
			authServiceMockBehavior: func(mockClient *mock_proto.MockAuthorizationClient,
				req *proto.Empty, expectedResponse *proto.IsAuthResponse) {
				mockClient.EXPECT().IsAuth(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "1"), "user_id", entity.UserID(1)),
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

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/is_auth", bytes.NewBuffer([]byte{}))
			r = r.WithContext(testCase.context)

			h.IsAuth(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	type AuthServiceMockBehavior func(
		mockClient *grpcAuth.MockAuthorizationClient,
		req *proto.LogoutRequest,
		expectedResponse *proto.LogoutResponse,
	)

	testTable := []struct {
		name               string
		expectedStatusCode int
		expectedJSON       string

		expectedRequest         *proto.LogoutRequest
		expectedResponse        *proto.LogoutResponse
		authServiceMockBehavior AuthServiceMockBehavior

		expectCookie bool

		context context.Context
	}{
		{
			name:               "OK test case 1",
			expectedStatusCode: 204,
			expectedJSON:       ``,
			expectedRequest: &proto.LogoutRequest{
				SessionToken: "ses1",
			},
			expectedResponse: &proto.LogoutResponse{},
			authServiceMockBehavior: func(mockClient *mock_proto.MockAuthorizationClient,
				req *proto.LogoutRequest, expectedResponse *proto.LogoutResponse) {
				mockClient.EXPECT().Logout(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context:      context.WithValue(context.Background(), "request_id", "1"),
			expectCookie: true,
		},
		{
			name:               "OK test case 1",
			expectedStatusCode: 204,
			expectedJSON:       ``,
			expectedRequest: &proto.LogoutRequest{
				SessionToken: "ses1",
			},
			expectedResponse: &proto.LogoutResponse{},
			authServiceMockBehavior: func(mockClient *mock_proto.MockAuthorizationClient,
				req *proto.LogoutRequest, expectedResponse *proto.LogoutResponse) {
				mockClient.EXPECT().Logout(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context:      context.WithValue(context.Background(), "request_id", "1"),
			expectCookie: false,
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

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/logout", bytes.NewBuffer([]byte{}))
			r = r.WithContext(testCase.context)
			if testCase.expectCookie {
				r.AddCookie(&http.Cookie{Name: "session_token", Value: "ses1"})
			}

			h.Logout(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}

func TestRegister(t *testing.T) {
	curTime := time.Time{}
	timeString := curTime.Format(time.RFC3339Nano)

	type AuthServiceMockBehavior func(
		mockClient *grpcAuth.MockAuthorizationClient,
		req *proto.LoginUserRequest,
		expectedResponse *proto.LoginUserResponse,
	)
	type mockArgs struct {
		User entity.User
	}
	type mockReturn struct {
		User         entity.User
		RegisterErrs []errs.ErrorInfo
		GetUserErr   errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	type test struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		Request          []byte
		ExpectedResponse expectedResponse
		Ctx              context.Context

		expectedRequest         *proto.LoginUserRequest
		expectedResponse        *proto.LoginUserResponse
		authServiceMockBehavior AuthServiceMockBehavior
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				User: users[0],
			},
			MockReturn: mockReturn{
				User: users[0],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","nickname":"%s","password":"%s"}`,
				users[0].Email, users[0].Nickname, users[0].Password)),
			expectedRequest: &proto.LoginUserRequest{
				UserId:     int64(users[0].UserID),
				Email:      users[0].Email,
				Nickname:   users[0].Nickname,
				Password:   users[0].Password,
				AvatarURL:  users[0].AvatarURL,
				RegisterAt: timeString,
			},
			expectedResponse: &proto.LoginUserResponse{
				UserId:     int64(users[0].UserID),
				Email:      users[0].Email,
				Nickname:   users[0].Nickname,
				Password:   users[0].Password,
				AvatarURL:  users[0].AvatarURL,
				Valid:      true,
				LocalError: 0,
				ExpiresAt:  timeString,
			},
			authServiceMockBehavior: func(mockClient *grpcAuth.MockAuthorizationClient,
				req *proto.LoginUserRequest, expectedResponse *proto.LoginUserResponse) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			ExpectedResponse: expectedResponse{
				Body: MakeUserResponseBodyWithBounds(users[0]),
				Code: 200,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				User: users[0],
			},
			MockReturn: mockReturn{
				User: users[0],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","nickname":"%s","password":"%s"}`,
				users[0].Email, users[0].Nickname, users[0].Password)),
			expectedRequest: &proto.LoginUserRequest{
				UserId:     int64(users[0].UserID),
				Email:      users[0].Email,
				Nickname:   users[0].Nickname,
				Password:   users[0].Password,
				AvatarURL:  users[0].AvatarURL,
				RegisterAt: timeString,
			},
			expectedResponse: &proto.LoginUserResponse{
				UserId:     int64(users[0].UserID),
				Email:      users[0].Email,
				Nickname:   users[0].Nickname,
				Password:   users[0].Password,
				AvatarURL:  users[0].AvatarURL,
				Valid:      false,
				LocalError: 11,
				ExpiresAt:  timeString,
			},
			authServiceMockBehavior: func(mockClient *grpcAuth.MockAuthorizationClient,
				req *proto.LoginUserRequest, expectedResponse *proto.LoginUserResponse) {
				mockClient.EXPECT().Login(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	for _, curTest := range tests {
		ctrl := gomock.NewController(t)

		mockAuthServiceClient := mock_proto.NewMockAuthorizationClient(ctrl)
		curTest.authServiceMockBehavior(mockAuthServiceClient, curTest.expectedRequest, curTest.expectedResponse)
		serviceMock := mock_service.NewMockIService(ctrl)

		h := handler.NewAPIHandler(serviceMock, zap.L(), nil, mockAuthServiceClient, nil, nil)
		r := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(curTest.Request))
		r = r.WithContext(curTest.Ctx)
		w := httptest.NewRecorder()
		serviceMock.EXPECT().RegisterUser(curTest.Ctx, gomock.Any()).
			Return(curTest.MockReturn.RegisterErrs).MaxTimes(1)
		serviceMock.EXPECT().GetUserByEmail(curTest.Ctx, curTest.MockArgs.User.Email).
			Return(curTest.MockReturn.User, curTest.MockReturn.GetUserErr).MaxTimes(1)
		h.Register(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}
