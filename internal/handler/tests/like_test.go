package tests

import (
	"bytes"
	"context"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	mock_service "harmonica/mocks/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"harmonica/internal/microservices/like/proto"
	grpcLike "harmonica/mocks/microservices/like/proto"
	mock_proto "harmonica/mocks/microservices/like/proto"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler_SetLike(t *testing.T) {
	type LikeServiceMockBehavior func(
		mockClient *grpcLike.MockLikeClient,
		req *proto.MakeLikeRequest,
		expectedResponse *proto.MakeLikeResponse,
	)
	testTable := []struct {
		name               string
		expectedStatusCode int
		expectedJSON       string

		expectedRequest         *proto.MakeLikeRequest
		expectedResponse        *proto.MakeLikeResponse
		likeServiceMockBehavior LikeServiceMockBehavior

		context context.Context
	}{
		{
			name:               "OK test 1",
			expectedStatusCode: 200,
			expectedJSON:       "null",
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      true,
				LocalError: 0,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().SetLike(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 1",
			expectedStatusCode: 400,
			expectedJSON:       MakeErrorResponse(errs.ErrInvalidSlug),
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  -1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      true,
				LocalError: 0,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().SetLike(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 2",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrDBInternal),
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      false,
				LocalError: 11,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().SetLike(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 3",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrGRPCWentWrong),
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      false,
				LocalError: 11,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().SetLike(gomock.Any(), req).Return(expectedResponse, errs.ErrGRPCWentWrong).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLikeServiceClient := mock_proto.NewMockLikeClient(ctrl)
			testCase.likeServiceMockBehavior(mockLikeServiceClient, testCase.expectedRequest, testCase.expectedResponse)
			serviceMock := mock_service.NewMockIService(ctrl)

			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, mockLikeServiceClient)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/pins/{pin_id}/like", bytes.NewBuffer([]byte{}))
			r.SetPathValue("pin_id", fmt.Sprintf(`%d`, testCase.expectedRequest.PinId))

			r = r.WithContext(testCase.context)

			h.CreateLike(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}

func TestHandler_DeleteLike(t *testing.T) {
	type LikeServiceMockBehavior func(
		mockClient *grpcLike.MockLikeClient,
		req *proto.MakeLikeRequest,
		expectedResponse *proto.MakeLikeResponse,
	)
	testTable := []struct {
		name               string
		expectedStatusCode int
		expectedJSON       string

		expectedRequest         *proto.MakeLikeRequest
		expectedResponse        *proto.MakeLikeResponse
		likeServiceMockBehavior LikeServiceMockBehavior

		context context.Context
	}{
		{
			name:               "OK test 1",
			expectedStatusCode: 200,
			expectedJSON:       "null",
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      true,
				LocalError: 0,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().ClearLike(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 1",
			expectedStatusCode: 400,
			expectedJSON:       MakeErrorResponse(errs.ErrInvalidSlug),
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  -1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      true,
				LocalError: 0,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().ClearLike(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 2",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrDBInternal),
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      false,
				LocalError: 11,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().ClearLike(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 3",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrGRPCWentWrong),
			expectedRequest: &proto.MakeLikeRequest{
				PinId:  1,
				UserId: 1,
			},
			expectedResponse: &proto.MakeLikeResponse{
				Valid:      false,
				LocalError: 11,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.MakeLikeRequest, expectedResponse *proto.MakeLikeResponse) {
				mockClient.EXPECT().ClearLike(gomock.Any(), req).Return(expectedResponse, errs.ErrGRPCWentWrong).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(), "request_id", "req_id"), "user_id", entity.UserID(1)),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLikeServiceClient := mock_proto.NewMockLikeClient(ctrl)
			testCase.likeServiceMockBehavior(mockLikeServiceClient, testCase.expectedRequest, testCase.expectedResponse)
			serviceMock := mock_service.NewMockIService(ctrl)

			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, mockLikeServiceClient)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/pins/{pin_id}/like", bytes.NewBuffer([]byte{}))
			r.SetPathValue("pin_id", fmt.Sprintf(`%d`, testCase.expectedRequest.PinId))

			r = r.WithContext(testCase.context)

			h.DeleteLike(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}

func TestHandler_UsersLike(t *testing.T) {
	type LikeServiceMockBehavior func(
		mockClient *grpcLike.MockLikeClient,
		req *proto.GetUsersLikedRequest,
		expectedResponse *proto.GetUsersLikedResponse,
	)
	testTable := []struct {
		name               string
		expectedStatusCode int
		expectedJSON       string

		expectedRequest         *proto.GetUsersLikedRequest
		expectedResponse        *proto.GetUsersLikedResponse
		likeServiceMockBehavior LikeServiceMockBehavior

		context context.Context
	}{
		{
			name:               "OK test 1",
			expectedStatusCode: 200,
			expectedJSON:       `{"users":null}`,

			expectedRequest: &proto.GetUsersLikedRequest{
				PinId: 1,
				Limit: 20,
			},
			expectedResponse: &proto.GetUsersLikedResponse{
				Valid:      true,
				LocalError: 0,
				Users:      []*proto.UserResponse{},
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.GetUsersLikedRequest, expectedResponse *proto.GetUsersLikedResponse) {
				mockClient.EXPECT().GetUsersLiked(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			name:               "Error test 1",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrGRPCWentWrong),

			expectedRequest: &proto.GetUsersLikedRequest{
				PinId: 1,
				Limit: 20,
			},
			expectedResponse: &proto.GetUsersLikedResponse{
				Valid:      true,
				LocalError: 0,
				Users:      []*proto.UserResponse{},
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.GetUsersLikedRequest, expectedResponse *proto.GetUsersLikedResponse) {
				mockClient.EXPECT().GetUsersLiked(gomock.Any(), req).Return(expectedResponse, errs.ErrGRPCWentWrong)
			},
			context: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			name:               "Error test 2",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrDBInternal),

			expectedRequest: &proto.GetUsersLikedRequest{
				PinId: 1,
				Limit: 20,
			},
			expectedResponse: &proto.GetUsersLikedResponse{
				Valid:      false,
				LocalError: 11,
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.GetUsersLikedRequest, expectedResponse *proto.GetUsersLikedResponse) {
				mockClient.EXPECT().GetUsersLiked(gomock.Any(), req).Return(expectedResponse, nil)
			},
			context: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLikeServiceClient := mock_proto.NewMockLikeClient(ctrl)
			testCase.likeServiceMockBehavior(mockLikeServiceClient, testCase.expectedRequest, testCase.expectedResponse)
			serviceMock := mock_service.NewMockIService(ctrl)

			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, mockLikeServiceClient)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/likes/{pin_id}/users", bytes.NewBuffer([]byte{}))
			r.SetPathValue("pin_id", fmt.Sprintf(`%d`, testCase.expectedRequest.PinId))

			r = r.WithContext(testCase.context)

			h.UsersLiked(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}

func TestHandler_GetFavorites(t *testing.T) {
	type LikeServiceMockBehavior func(
		mockClient *grpcLike.MockLikeClient,
		req *proto.GetFavoritesRequest,
		expectedResponse *proto.GetFavoritesResponse,
	)
	testTable := []struct {
		name               string
		expectedStatusCode int
		expectedJSON       string

		expectedRequest         *proto.GetFavoritesRequest
		expectedResponse        *proto.GetFavoritesResponse
		likeServiceMockBehavior LikeServiceMockBehavior

		context context.Context
	}{
		{
			name:               "OK test 1",
			expectedStatusCode: 200,
			expectedJSON:       `{"pins":null}`,
			expectedRequest: &proto.GetFavoritesRequest{
				UserId: 1,
				Limit:  40,
				Offset: 0,
			},
			expectedResponse: &proto.GetFavoritesResponse{
				Valid:      true,
				LocalError: 0,
				Pins:       []*proto.FeedPin{},
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.GetFavoritesRequest, expectedResponse *proto.GetFavoritesResponse) {
				mockClient.EXPECT().GetFavorites(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(),
				"request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 1",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrDBInternal),
			expectedRequest: &proto.GetFavoritesRequest{
				UserId: 1,
				Limit:  40,
				Offset: 0,
			},
			expectedResponse: &proto.GetFavoritesResponse{
				Valid:      false,
				LocalError: 11,
				Pins:       []*proto.FeedPin{},
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.GetFavoritesRequest, expectedResponse *proto.GetFavoritesResponse) {
				mockClient.EXPECT().GetFavorites(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(),
				"request_id", "req_id"), "user_id", entity.UserID(1)),
		},
		{
			name:               "Error test 2",
			expectedStatusCode: 500,
			expectedJSON:       MakeErrorResponse(errs.ErrGRPCWentWrong),
			expectedRequest: &proto.GetFavoritesRequest{
				UserId: 1,
				Limit:  40,
				Offset: 0,
			},
			expectedResponse: &proto.GetFavoritesResponse{
				Valid:      true,
				LocalError: 0,
				Pins:       []*proto.FeedPin{},
			},
			likeServiceMockBehavior: func(mockClient *mock_proto.MockLikeClient,
				req *proto.GetFavoritesRequest, expectedResponse *proto.GetFavoritesResponse) {
				mockClient.EXPECT().GetFavorites(gomock.Any(), req).Return(expectedResponse, errs.ErrGRPCWentWrong).AnyTimes()
			},
			context: context.WithValue(context.WithValue(context.Background(),
				"request_id", "req_id"), "user_id", entity.UserID(1)),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLikeServiceClient := mock_proto.NewMockLikeClient(ctrl)
			testCase.likeServiceMockBehavior(mockLikeServiceClient, testCase.expectedRequest, testCase.expectedResponse)
			serviceMock := mock_service.NewMockIService(ctrl)

			h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, mockLikeServiceClient)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/favorites", bytes.NewBuffer([]byte{}))

			r = r.WithContext(testCase.context)

			h.GetFavorites(w, r)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedJSON, w.Body.String())
		})
	}
}

/*

func TestCreateLike(t *testing.T) {
	type mockArgs struct {
		Ctx    context.Context
		PinId  entity.PinID
		UserId entity.UserID
		Slug   any // pinId
		Times  int
	}
	type mockReturn struct {
		Err errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	type test struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Ctx:    context.WithValue(context.Background(), "user_id", entity.UserID(0)),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   1,
				Times:  1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{},
			},
			ExpectedResponse: expectedResponse{
				Body: "null",
				Code: 200,
			},
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				Ctx:    context.WithValue(context.Background(), "user_id", entity.UserID(0)),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   1,
				Times:  1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
		},
		{
			Name: "Uncorrect test 2",
			MockArgs: mockArgs{
				Ctx:    context.WithValue(context.Background(), "user_id", entity.UserID(0)),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   "Abc",
				Times:  0,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
		},
	}

	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L())
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/pins/", nil)
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.MockArgs.Slug))
		ctx := context.WithValue(curTest.MockArgs.Ctx, "request_id", "req_id")
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		serviceMock.EXPECT().SetLike(ctx, curTest.MockArgs.PinId, curTest.MockArgs.UserId).
			Return(curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.CreateLike(w, r)
		assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
		assert.Equal(t, curTest.ExpectedResponse.Body, w.Body.String())
	}
}

func TestDeleteLike(t *testing.T) {
	type mockArgs struct {
		Ctx    context.Context
		PinId  entity.PinID
		UserId entity.UserID
		Slug   any // pinId
		Times  int
	}
	type mockReturn struct {
		Err errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	type test struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Ctx:    context.WithValue(context.Background(), "user_id", entity.UserID(0)),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   1,
				Times:  1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{},
			},
			ExpectedResponse: expectedResponse{
				Body: "null",
				Code: 200,
			},
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				Ctx:    context.WithValue(context.Background(), "user_id", entity.UserID(0)),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   1,
				Times:  1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
		},
		{
			Name: "Uncorrect test 2",
			MockArgs: mockArgs{
				Ctx:    context.WithValue(context.Background(), "user_id", entity.UserID(0)),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   "Abc",
				Times:  0,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
		},
	}

	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L())
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/pins/", nil)
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.MockArgs.Slug))
		ctx := context.WithValue(curTest.MockArgs.Ctx, "request_id", "req_id")
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		serviceMock.EXPECT().ClearLike(ctx, curTest.MockArgs.PinId, curTest.MockArgs.UserId).
			Return(curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.DeleteLike(w, r)
		assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
		assert.Equal(t, curTest.ExpectedResponse.Body, w.Body.String())
	}
}

func TestUsersLiked(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		PinId entity.PinID
		Limit int
		Slug  any // pinId
		Times int
	}
	type mockReturn struct {
		List entity.UserList
		Err  errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	type test struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		ExpectedResponse expectedResponse
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				PinId: entity.PinID(1),
				Slug:  1,
				Limit: 20,
				Times: 1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{},
			},
			ExpectedResponse: expectedResponse{
				Body: `{"users":null}`,
				Code: 200,
			},
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				PinId: entity.PinID(1),
				Slug:  1,
				Limit: 20,
				Times: 1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
		},
		{
			Name: "Uncorrect test 2",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				PinId: entity.PinID(1),
				Slug:  "Abc",
				Limit: 20,
				Times: 0,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
		},
	}

	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L())
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/pins/", nil)
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.MockArgs.Slug))
		ctx := context.WithValue(curTest.MockArgs.Ctx, "request_id", "req_id")
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		serviceMock.EXPECT().GetUsersLiked(ctx, curTest.MockArgs.PinId, curTest.MockArgs.Limit).
			Return(curTest.MockReturn.List, curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.UsersLiked(w, r)
		assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
		assert.Equal(t, curTest.ExpectedResponse.Body, w.Body.String())
	}
}
*/
