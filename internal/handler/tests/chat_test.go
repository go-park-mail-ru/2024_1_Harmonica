package tests

/*
import (
	"bytes"
	"context"
	"encoding/json"
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
)

func TestAPIHandler_SendMessage(t *testing.T) {
	type mockArgs struct {
		Message entity.Message
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
		UserId           entity.UserID
		ExpectedResponse expectedResponse
	}{
		{
			Name: "OK test 1",
			MockArgs: mockArgs{
				Message: entity.Message{
					ReceiverId: 1,
					SenderId:   2,
					Text:       "Hello!",
				},
			},
			MockReturn: mockReturn{},
			Slug:       1,
			UserId:     2,
			ExpectedResponse: expectedResponse{
				Body: `null`,
				Code: http.StatusOK,
			},
		},
		{
			Name:       "Error test 1",
			MockArgs:   mockArgs{},
			MockReturn: mockReturn{},
			Slug:       "abc",
			UserId:     2,
			ExpectedResponse: expectedResponse{
				Body: `{"code":12,"message":"invalid slug parameter"}`,
				Code: http.StatusBadRequest,
			},
		},
		{
			Name: "Error test 2",
			MockArgs: mockArgs{
				Message: entity.Message{
					Text: "   ",
				},
			},
			MockReturn: mockReturn{},
			Slug:       1,
			UserId:     2,
			ExpectedResponse: expectedResponse{
				Body: `{"code":5,"message":"validation conditions are not met"}`,
				Code: http.StatusBadRequest,
			},
		},
		{
			Name: "Error test 3",
			MockArgs: mockArgs{
				Message: entity.Message{
					ReceiverId: 1,
					SenderId:   2,
					Text:       "Hello!",
				},
			},
			MockReturn: mockReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
			},
			Slug:   1,
			UserId: 2,
			ExpectedResponse: expectedResponse{
				Body: `{"code":11,"message":"internal db error"}`,
				Code: http.StatusInternalServerError,
			},
		},
	}

	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)
	for _, curTest := range tests {
		t.Run(curTest.Name, func(t *testing.T) {
			reqBytes, err := json.Marshal(curTest.MockArgs.Message)
			if err != nil {
				t.Errorf("error marshaling request body: %v", err)
			}
			r := httptest.NewRequest(http.MethodPost, "/api/v1/messages/{receiver_id}", bytes.NewBuffer(reqBytes))
			w := httptest.NewRecorder()
			r.SetPathValue("receiver_id", fmt.Sprintf(`%d`, curTest.Slug))
			ctx := context.WithValue(context.Background(), "user_id", curTest.UserId)
			r = r.WithContext(ctx)
			serviceMock.EXPECT().CreateMessage(ctx, curTest.MockArgs.Message).Return(curTest.MockReturn.ErrorInfo).MaxTimes(1)
			h.SendMessage(w, r)
			assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
			assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
		})
	}
}

func TestAPIHandler_ReadMessages(t *testing.T) {
	type mockArgs struct {
		UserId1 entity.UserID
		UserId2 entity.UserID
	}
	type mockReturn struct {
		Messages  entity.Messages
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
			Name: "OK test 1",
			MockArgs: mockArgs{
				UserId1: 1,
				UserId2: 2,
			},
			MockReturn: mockReturn{
				Messages: entity.Messages{Messages: []entity.MessageResponse{{Text: "hello"}}},
			},
			Slug: 1,
			ExpectedResponse: expectedResponse{
				Body: `{"messages":[{"sender_id":0,"receiver_id":0,"text":"hello","status":"","sent_at":"0001-01-01T00:00:00Z"}]}`,
				Code: http.StatusOK,
			},
		},
		{
			Name: "Error test 1",
			MockArgs: mockArgs{
				UserId1: 1,
				UserId2: 2,
			},
			MockReturn: mockReturn{},
			Slug:       "abc",
			ExpectedResponse: expectedResponse{
				Body: `{"code":12,"message":"invalid slug parameter"}`,
				Code: http.StatusBadRequest,
			},
		},
	}

	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)
	for _, curTest := range tests {
		t.Run(curTest.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/v1/messages/{user_id}", nil)
			w := httptest.NewRecorder()
			r.SetPathValue("user_id", fmt.Sprintf(`%d`, curTest.Slug))
			ctx := context.WithValue(context.Background(), "user_id", curTest.MockArgs.UserId2)
			r = r.WithContext(ctx)
			serviceMock.EXPECT().GetMessages(ctx, curTest.MockArgs.UserId1, curTest.MockArgs.UserId2).
				Return(curTest.MockReturn.Messages, curTest.MockReturn.ErrorInfo).MaxTimes(1)
			h.ReadMessages(w, r)
			assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
			assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
		})
	}
}
*/
