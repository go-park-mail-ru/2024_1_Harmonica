package tests

import (
	"bytes"
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
)

func TestAPIHandler_UpdateDraft(t *testing.T) {
	type mockReturn struct {
		ErrorInfo errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	draft := entity.Draft{SenderId: 2, ReceiverId: 123, Text: "Test draft"}
	tests := []struct {
		Name             string
		Slug             any
		RequestJSON      string
		UserId           entity.UserID
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}{
		{
			Name:        "Error test - Invalid slug",
			Slug:        "abc",
			RequestJSON: `{"text":"Test draft"}`,
			UserId:      draft.SenderId,
			MockReturn:  mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: `{"code":12,"message":"invalid slug parameter"}`,
				Code: http.StatusBadRequest,
			},
		},
		{
			Name:        "Error test - Error reading request body",
			Slug:        draft.ReceiverId,
			RequestJSON: `invalid json`,
			UserId:      draft.SenderId,
			MockReturn:  mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: `{"code":4,"message":"error reading request body"}`,
				Code: http.StatusBadRequest,
			},
		},
		{
			Name:        "Error test - Internal DB error",
			Slug:        draft.ReceiverId,
			RequestJSON: `{"text":"Test draft"}`,
			UserId:      draft.SenderId,
			MockReturn: mockReturn{
				ErrorInfo: errs.ErrorInfo{GeneralErr: errs.ErrDBInternal, LocalErr: errs.ErrDBInternal},
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
			reqBody := bytes.NewBufferString(curTest.RequestJSON)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/drafts/{receiver_id}", reqBody)
			req.SetPathValue("receiver_id", fmt.Sprintf(`%d`, curTest.Slug))
			ctx := context.WithValue(req.Context(), "user_id", draft.SenderId)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()
			if curTest.ExpectedResponse.Code == 500 {
				serviceMock.EXPECT().UpdateDraft(ctx, draft).Return(curTest.MockReturn.ErrorInfo)
			}
			handlerFunc := http.HandlerFunc(h.UpdateDraft)
			handlerFunc.ServeHTTP(w, req)
			assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
			assert.JSONEq(t, curTest.ExpectedResponse.Body, w.Body.String())
		})
	}
}
