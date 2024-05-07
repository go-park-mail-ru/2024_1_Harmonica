package tests

import (
	"bytes"
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	mock_service "harmonica/mocks/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler_Search(t *testing.T) {
	type MainServiceMockBehavior func(
		mockClient *mock_service.MockIService,
		reqQuery string,
		expectedResponseRes entity.SearchResult,
		expectedResponseErr errs.ErrorInfo,
	)
	testTable := []struct {
		name               string
		expectedStatusCode int
		expectedJSON       string

		reqQuery                string
		expectedResponseRes     entity.SearchResult
		expectedResponseErr     errs.ErrorInfo
		mainServiceMockBehavior MainServiceMockBehavior

		context context.Context
	}{
		{
			name:                "OK test 1",
			expectedStatusCode:  200,
			expectedJSON:        `{"users":null,"pins":null,"boards":null}`,
			reqQuery:            "miracle",
			expectedResponseRes: entity.SearchResult{},
			expectedResponseErr: errs.ErrorInfo{},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqQuery string,
				expectedResponseRes entity.SearchResult, expectedResponseErr errs.ErrorInfo) {
				mockClient.EXPECT().Search(gomock.Any(), reqQuery).Return(expectedResponseRes, expectedResponseErr).AnyTimes()
			},
			context: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			name:                "Error test 1",
			expectedStatusCode:  500,
			expectedJSON:        MakeErrorResponse(errs.ErrDBInternal),
			reqQuery:            "miracle",
			expectedResponseRes: entity.SearchResult{},
			expectedResponseErr: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqQuery string,
				expectedResponseRes entity.SearchResult, expectedResponseErr errs.ErrorInfo) {
				mockClient.EXPECT().Search(gomock.Any(), reqQuery).Return(expectedResponseRes, expectedResponseErr).AnyTimes()
			},
			context: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			name:                "Error test 2",
			expectedStatusCode:  400,
			expectedJSON:        MakeErrorResponse(errs.ErrInvalidSlug),
			reqQuery:            "",
			expectedResponseRes: entity.SearchResult{},
			expectedResponseErr: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			mainServiceMockBehavior: func(mockClient *mock_service.MockIService, reqQuery string,
				expectedResponseRes entity.SearchResult, expectedResponseErr errs.ErrorInfo) {
				mockClient.EXPECT().Search(gomock.Any(), reqQuery).Return(expectedResponseRes, expectedResponseErr).AnyTimes()
			},
			context: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceMock := mock_service.NewMockIService(ctrl)
			testCase.mainServiceMockBehavior(serviceMock, testCase.reqQuery,
				testCase.expectedResponseRes, testCase.expectedResponseErr)
			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/search/{search_query}", bytes.NewBuffer([]byte{}))
			r.SetPathValue("search_query", testCase.reqQuery)

			r = r.WithContext(testCase.context)

			h.Search(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}
