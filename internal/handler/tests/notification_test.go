package tests

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	mock_service "harmonica/mocks/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIHandler_GetUnreadNotifications(t *testing.T) {
	type mockArgs struct {
		UserId entity.UserID
	}
	type mockReturn struct {
		Notifications entity.Notifications
		ErrorInfo     errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	tests := []struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}{
		{
			Name: "OK test 1",
			MockArgs: mockArgs{
				UserId: 1,
			},
			MockReturn: mockReturn{
				Notifications: entity.Notifications{
					Notifications: []entity.NotificationResponse{
						{
							NotificationId: 1,
							UserId:         1,
							Type:           "message",
							TriggeredByUser: entity.TriggeredByUser{
								UserId:   2,
								Nickname: "user2",
							},
							Message: entity.MessageNotificationResponse{
								Text: "You have a new message",
							},
							CreatedAt: time.Date(2024, 5, 28, 19, 27, 40, 0, time.UTC),
						},
					},
				},
				ErrorInfo: errs.ErrorInfo{},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"notifications":[{"notification_id":1,"user_id":1,"type":"message","triggered_by_user":{"user_id":2,"nickname":"user2","avatar_url":""},"pin":{"pin_id":0,"content_url":""},"comment":{"comment_id":0,"text":""},"message":{"text":"You have a new message"},"created_at":"2024-05-28T19:27:40Z"}]}`,
				Code: http.StatusOK,
			},
		},
		{
			Name: "Error test 1",
			MockArgs: mockArgs{
				UserId: 1,
			},
			MockReturn: mockReturn{
				ErrorInfo: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"code":11,"message":"internal db error"}`,
				Code: http.StatusInternalServerError,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)

	for _, curTest := range tests {
		t.Run(curTest.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(r.Context(), "user_id", curTest.MockArgs.UserId)
			r = r.WithContext(ctx)

			serviceMock.EXPECT().GetUnreadNotifications(ctx, curTest.MockArgs.UserId).
				Return(curTest.MockReturn.Notifications, curTest.MockReturn.ErrorInfo).MaxTimes(1)

			h.GetUnreadNotifications(w, r)

			assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
			assert.JSONEq(t, curTest.ExpectedResponse.Body, w.Body.String())
		})
	}
}

func TestAPIHandler_ReadNotification(t *testing.T) {
	type mockArgs struct {
		NotificationID entity.NotificationID
		UserId         entity.UserID
	}
	type mockReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	tests := []struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		Slug             any
		ExpectedResponse expectedResponse
	}{
		{
			Name: "OK test",
			MockArgs: mockArgs{
				NotificationID: 1,
				UserId:         2,
			},
			MockReturn: mockReturn{},
			Slug:       1,
			ExpectedResponse: expectedResponse{
				Body: `null`,
				Code: http.StatusOK,
			},
		},
		{
			Name: "Error test - Invalid notification ID",
			MockArgs: mockArgs{
				NotificationID: 0, // Invalid notification ID
				UserId:         2,
			},
			MockReturn: mockReturn{},
			Slug:       "abc",
			ExpectedResponse: expectedResponse{
				Body: `{"code":12,"message":"invalid slug parameter"}`,
				Code: http.StatusBadRequest,
			},
		},
		{
			Name: "Error test - Service error",
			MockArgs: mockArgs{
				NotificationID: 1,
				UserId:         2,
			},
			MockReturn: mockReturn{
				ErrorInfo: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			Slug: 1,
			ExpectedResponse: expectedResponse{
				Body: `{"code":11,"message":"internal db error"}`,
				Code: http.StatusInternalServerError,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)

	for _, curTest := range tests {
		t.Run(curTest.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/read/{notification_id}", nil)
			w := httptest.NewRecorder()
			r.SetPathValue("notification_id", fmt.Sprintf(`%d`, curTest.Slug))
			ctx := context.WithValue(r.Context(), "user_id", curTest.MockArgs.UserId)
			r = r.WithContext(ctx)
			serviceMock.EXPECT().ReadNotification(ctx, curTest.MockArgs.NotificationID, curTest.MockArgs.UserId).
				Return(curTest.MockReturn.ErrorInfo).MaxTimes(1)
			h.ReadNotification(w, r)
			assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
			assert.JSONEq(t, curTest.ExpectedResponse.Body, w.Body.String())
		})
	}
}

func TestAPIHandler_ReadAllNotifications(t *testing.T) {
	type mockArgs struct {
		UserID entity.UserID
	}
	type mockReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	tests := []struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}{
		{
			Name: "OK test",
			MockArgs: mockArgs{
				UserID: 1,
			},
			MockReturn: mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: `null`,
				Code: http.StatusOK,
			},
		},
		{
			Name: "Error test - Service error",
			MockArgs: mockArgs{
				UserID: 1,
			},
			MockReturn: mockReturn{
				ErrorInfo: errs.ErrorInfo{
					GeneralErr: errs.ErrDBInternal,
					LocalErr:   errs.ErrDBInternal,
				},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"code":11,"message":"internal db error"}`,
				Code: http.StatusInternalServerError,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)

	for _, curTest := range tests {
		t.Run(curTest.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/read/all", nil)
			w := httptest.NewRecorder()
			ctx := context.WithValue(r.Context(), "user_id", curTest.MockArgs.UserID)
			r = r.WithContext(ctx)
			serviceMock.EXPECT().ReadAllNotifications(ctx, curTest.MockArgs.UserID).
				Return(curTest.MockReturn.ErrorInfo).MaxTimes(1)
			h.ReadAllNotifications(w, r)
			assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
			assert.JSONEq(t, curTest.ExpectedResponse.Body, w.Body.String())
		})
	}
}
