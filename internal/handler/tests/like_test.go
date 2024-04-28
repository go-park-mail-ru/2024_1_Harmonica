package tests
/*
import (
	"context"
	"fmt"
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