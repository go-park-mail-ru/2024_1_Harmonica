package tests

import (
	"bytes"
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	mock_proto "harmonica/mocks/microservices/image/proto"
	mock_service "harmonica/mocks/service"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var ContextWithUserID = context.WithValue(context.Background(), "user_id", entity.UserID(1))
var ContextWithUserAndRequestID = context.WithValue(ContextWithUserID, "request_id", "test_req_id")

func TestHandler_GetUser(t *testing.T) {
	type MainServiceMockBehavior func(
		mockClient *mock_service.MockIService,
		reqNickname string,
		reqUserID entity.UserID,
		expectedRes entity.UserProfileResponse,
		expectedErr errs.ErrorInfo,
	)

	testTable := []struct {
		name string

		expectedStatusCode int
		expectedJSON       string

		reqNickname             string
		reqUserID               entity.UserID
		resProfile              entity.UserProfileResponse
		resErr                  errs.ErrorInfo
		mainServiceMockBehavior MainServiceMockBehavior

		context context.Context
	}{
		{
			name:               "OK test 1",
			expectedStatusCode: 200,
			expectedJSON: `{"user":{"user_id":1,"email":"e@e.ru","nickname":"MIRACLE","avatar_url":"",` +
				`"avatar_width":0,"avatar_height":0},"subscriptions_count":0,"subscribers_count":0,"is_subscribed` +
				`":false,"is_owner":false}`,
			reqNickname: "MIRACLE",
			reqUserID:   entity.UserID(1),
			resProfile: entity.UserProfileResponse{
				User: entity.UserResponse{
					UserId:   entity.UserID(1),
					Nickname: "MIRACLE",
					Email:    "e@e.ru",
				},
			},
			resErr: errs.ErrorInfo{},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqNickname string,
				reqUserID entity.UserID, expectedRes entity.UserProfileResponse, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().GetUserProfileByNickname(gomock.Any(), reqNickname, reqUserID).
					Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
		{
			name:               "Error test 1",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrDBInternal),
			reqNickname:        "MIRACLE",
			reqUserID:          entity.UserID(1),
			resProfile: entity.UserProfileResponse{
				User: entity.UserResponse{
					UserId:   entity.UserID(1),
					Nickname: "MIRACLE",
					Email:    "e@e.ru",
				},
			},
			resErr: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqNickname string,
				reqUserID entity.UserID, expectedRes entity.UserProfileResponse, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().GetUserProfileByNickname(gomock.Any(), reqNickname, reqUserID).
					Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
		{
			name:               "Error test 2",
			expectedStatusCode: 404,
			expectedJSON:       MakeErrorResponse(errs.ErrUserNotExist),
			reqNickname:        "MIRACLE",
			reqUserID:          entity.UserID(1),
			resProfile: entity.UserProfileResponse{
				User: entity.UserResponse{},
			},
			resErr: errs.ErrorInfo{},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqNickname string,
				reqUserID entity.UserID, expectedRes entity.UserProfileResponse, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().GetUserProfileByNickname(gomock.Any(), reqNickname, reqUserID).
					Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock_service.NewMockIService(ctrl)
			testCase.mainServiceMockBehavior(serviceMock, testCase.reqNickname, testCase.reqUserID,
				testCase.resProfile, testCase.resErr)
			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/users/{nickname}", bytes.NewBuffer([]byte{}))
			r = r.WithContext(testCase.context)
			r.SetPathValue("nickname", testCase.reqNickname)

			h.GetUser(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	type MainServiceMockBehavior func(
		mockClient *mock_service.MockIService,
		reqUser entity.User,
		reqContext context.Context,
		expectedRes entity.User,
		expectedErr errs.ErrorInfo,
	)

	testTable := []struct {
		name               string
		requestBody        url.Values
		expectedStatusCode int
		expectedJSON       string

		pathUserID string

		reqUser                 entity.User
		reqContext              context.Context
		expectedRes             entity.User
		expectedErr             errs.ErrorInfo
		mainServiceMockBehavior MainServiceMockBehavior

		context context.Context
	}{
		{
			name: "OK test 1",
			requestBody: url.Values{
				"user": []string{`{"user_id":1,"nickname":"marci"}`},
			},
			expectedStatusCode: 200,
			expectedJSON:       `{"user_id":1,"email":"","nickname":"marci","avatar_url":"","avatar_width":0,"avatar_height":0}`,

			pathUserID: "1",

			reqUser: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			reqContext: ContextWithUserAndRequestID,
			expectedRes: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			expectedErr: errs.ErrorInfo{},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqUser entity.User,
				reqContext context.Context, expectedRes entity.User, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().UpdateUser(gomock.Any(), reqUser).Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
		{
			name: "Error test 1",
			requestBody: url.Values{
				"user": []string{`{"user_id":1,"nickname":"marci"}`},
			},
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrDBInternal),

			pathUserID: "1",

			reqUser: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			reqContext: ContextWithUserAndRequestID,
			expectedRes: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			expectedErr: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqUser entity.User,
				reqContext context.Context, expectedRes entity.User, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().UpdateUser(gomock.Any(), reqUser).Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
		{
			name: "Error test 2",
			requestBody: url.Values{
				"user": []string{`{"user_id":1,"nickname":"marci"}`},
			},
			expectedStatusCode: 400,
			expectedJSON:       MakeErrorResponse(errs.ErrDiffUserId),

			pathUserID: "2",

			reqUser: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			reqContext: ContextWithUserAndRequestID,
			expectedRes: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			expectedErr: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqUser entity.User,
				reqContext context.Context, expectedRes entity.User, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().UpdateUser(gomock.Any(), reqUser).Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
		{
			name: "Error test 3",
			requestBody: url.Values{
				"user": []string{`{"user_id:1,"nickname":"marci"}`},
			},
			expectedStatusCode: 400,
			expectedJSON:       MakeErrorResponse(errs.ErrReadingRequestBody),

			pathUserID: "1",

			reqUser: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			reqContext: ContextWithUserAndRequestID,
			expectedRes: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			expectedErr: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqUser entity.User,
				reqContext context.Context, expectedRes entity.User, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().UpdateUser(gomock.Any(), reqUser).Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
		{
			name: "Error test 4",
			requestBody: url.Values{
				"user": []string{`{"user_id":1,"nickname":"marci","password":"123"}`},
			},
			expectedStatusCode: 400,
			expectedJSON:       MakeErrorResponse(errs.ErrInvalidInputFormat),

			pathUserID: "1",

			reqUser: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			reqContext: ContextWithUserAndRequestID,
			expectedRes: entity.User{
				UserID:   entity.UserID(1),
				Nickname: "marci",
			},
			expectedErr: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqUser entity.User,
				reqContext context.Context, expectedRes entity.User, expectedErr errs.ErrorInfo) {
				mockClient.EXPECT().UpdateUser(gomock.Any(), reqUser).Return(expectedRes, expectedErr).AnyTimes()
			},
			context: ContextWithUserAndRequestID,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock_service.NewMockIService(ctrl)

			imageMock := mock_proto.NewMockImageClient(ctrl)
			testCase.mainServiceMockBehavior(serviceMock, testCase.reqUser, testCase.reqContext,
				testCase.expectedRes, testCase.expectedErr)
			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, imageMock, nil)

			dataReader := testCase.requestBody.Encode()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/users/{user_id}", bytes.NewBuffer([]byte(dataReader)))
			r = r.WithContext(testCase.context)
			r.SetPathValue("user_id", testCase.pathUserID)
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			h.UpdateUser(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}
