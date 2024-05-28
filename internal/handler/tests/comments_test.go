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
	"strings"
	"testing"
)

func TestAPIHandler_AddComment(t *testing.T) {
	type mockArgs struct {
		Value  string
		PinId  entity.PinID
		UserId entity.UserID
	}
	type mockReturn struct {
		Pin            entity.PinPageResponse
		CommentId      entity.CommentID
		ErrorInfo1     errs.ErrorInfo
		ErrorInfo2     errs.ErrorInfo
		ErrorInfo3     errs.ErrorInfo
		User           entity.User
		NotificationId entity.NotificationID
	}
	type expectedResponse struct {
		Body string
		Code int
	}

	tests := []struct {
		Name             string
		Slug             any
		ReqBody          string
		MockArgs         mockArgs
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}{
		//{
		//	Name: "OK test 1",
		//	Slug: 123,
		//	MockArgs: mockArgs{
		//		Value:  "Test comment",
		//		PinId:  123,
		//		UserId: 456,
		//	},
		//	MockReturn: mockReturn{
		//		Pin: entity.PinPageResponse{
		//			PinAuthor: entity.PinAuthor{
		//				UserId: 789,
		//			},
		//		},
		//		CommentId:      789,
		//		NotificationId: 1001,
		//	},
		//	ExpectedResponse: expectedResponse{
		//		Body: "null",
		//		Code: http.StatusBadRequest,
		//	},
		//},
		{
			Name: "Error test 1",
			Slug: "abc",
			MockArgs: mockArgs{
				Value:  "Test comment",
				UserId: 456,
			},
			MockReturn: mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: `{"code":12,"message":"invalid slug parameter"}`,
				Code: http.StatusBadRequest,
			},
		},
		{
			Name:       "Error test 2",
			Slug:       3,
			ReqBody:    `{"value":"Test`,
			MockArgs:   mockArgs{},
			MockReturn: mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: `{"code":4,"message":"error reading request body"}`,
				Code: http.StatusBadRequest,
			},
		},
		{
			Name:    "Error test 3",
			Slug:    0,
			ReqBody: `{"value":"Test comment"}`,
			MockArgs: mockArgs{
				Value:  "Test comment",
				UserId: 456,
			},
			MockReturn: mockReturn{
				ErrorInfo1: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"code":11,"message":"internal db error"}`,
				Code: http.StatusInternalServerError,
			},
		},
		{
			Name:    "Error test 4",
			Slug:    0,
			ReqBody: `{"value":"Test comment"}`,
			MockArgs: mockArgs{
				Value:  "Test comment",
				UserId: 456,
			},
			MockReturn: mockReturn{
				ErrorInfo2: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"code":11,"message":"internal db error"}`,
				Code: http.StatusInternalServerError,
			},
		},
		{
			Name:    "Error test 5",
			Slug:    0,
			ReqBody: `{"value":"Test comment"}`,
			MockArgs: mockArgs{
				Value:  "Test comment",
				UserId: 456,
			},
			MockReturn: mockReturn{
				ErrorInfo3: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
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
	logger := zap.NewNop()
	h := handler.NewAPIHandler(serviceMock, logger, nil, nil, nil, nil)

	for _, curTest := range tests {
		t.Run(curTest.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/v1/comments/{pin_id}", strings.NewReader(curTest.ReqBody))
			r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.Slug))
			w := httptest.NewRecorder()
			ctx := context.WithValue(r.Context(), "user_id", curTest.MockArgs.UserId)
			r = r.WithContext(ctx)
			if curTest.ExpectedResponse.Code == 500 {
				serviceMock.EXPECT().AddComment(gomock.Any(), "Test comment", curTest.MockArgs.PinId, curTest.MockArgs.UserId).
					Return(curTest.MockReturn.Pin, curTest.MockReturn.CommentId, curTest.MockReturn.ErrorInfo1)
			}
			if curTest.MockReturn.ErrorInfo2 != (errs.ErrorInfo{}) || curTest.MockReturn.ErrorInfo3 != (errs.ErrorInfo{}) {
				serviceMock.EXPECT().CreateNotification(gomock.Any(), gomock.Any()).Return(curTest.MockReturn.NotificationId, curTest.MockReturn.ErrorInfo2)
			}
			if curTest.MockReturn.ErrorInfo2 == (errs.ErrorInfo{}) && curTest.MockReturn.ErrorInfo3 != (errs.ErrorInfo{}) {
				serviceMock.EXPECT().GetUserById(gomock.Any(), curTest.MockArgs.UserId).Return(curTest.MockReturn.User, curTest.MockReturn.ErrorInfo3)
			}
			h.AddComment(w, r)
			assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
			assert.JSONEq(t, curTest.ExpectedResponse.Body, w.Body.String())
		})
	}
}

func TestAPIHandler_GetComments(t *testing.T) {
	type mockReturn struct {
		Comments entity.GetCommentsResponse
		ErrInfo  errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	pinId := entity.PinID(1)
	tests := []struct {
		Name             string
		Slug             any
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}{
		{
			Name: "OK test",
			Slug: pinId,
			MockReturn: mockReturn{
				Comments: entity.GetCommentsResponse{
					Comments: []entity.CommentResponse{
						{
							CommentId: 123,
							Value:     "Test comment 1",
							CommentAuthor: entity.CommentAuthor{
								UserId:    123,
								Nickname:  "user1",
								AvatarURL: "http://example.com/avatar1.jpg",
							},
						},
					},
				},
				ErrInfo: errs.ErrorInfo{},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"comments":[{"comment_id":123,"value":"Test comment 1","user":{"user_id":123,"nickname":"user1","avatar_url":"http://example.com/avatar1.jpg"}}]}`,
				Code: http.StatusOK,
			},
		},
		{
			Name: "Error test 1",
			Slug: "abc",
			MockReturn: mockReturn{
				ErrInfo: errs.ErrorInfo{
					LocalErr: errs.ErrInvalidSlug,
				},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"code":12,"message":"invalid slug parameter"}`,
				Code: http.StatusBadRequest,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	serviceMock := mock_service.NewMockIService(ctrl)
	logger := zap.NewNop()
	h := handler.NewAPIHandler(serviceMock, logger, nil, nil, nil, nil)

	for _, curTest := range tests {
		t.Run(curTest.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/v1/comments/{pin_id}", nil)
			r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.Slug))
			w := httptest.NewRecorder()
			ctx := context.WithValue(r.Context(), "user_id", 123)
			r = r.WithContext(ctx)

			serviceMock.EXPECT().GetComments(gomock.Any(), pinId).Return(curTest.MockReturn.Comments, curTest.MockReturn.ErrInfo).MaxTimes(1)

			handlerFunc := http.HandlerFunc(h.GetComments)
			handlerFunc.ServeHTTP(w, r)

			assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
			assert.JSONEq(t, curTest.ExpectedResponse.Body, w.Body.String())
		})
	}
}
