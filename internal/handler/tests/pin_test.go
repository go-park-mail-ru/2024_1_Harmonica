package tests

import (
	"bytes"
	"context"
	"encoding/json"
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

func MakeGetPinResponse(pin entity.PinPageResponse) string {
	return fmt.Sprintf(`{"pin_id":%d,"created_at":"0001-01-01T00:00:00Z","title":"%s","description":"%s",`+
		`"allow_comments":%t,"click_url":"%s","content_url":"%s","likes_count":%d,"is_owner":%t,"is_liked":%t,`+
		`"author":{"user_id":%d,"nickname":"%s","avatar_url":"%s"}}`,
		pin.PinId, pin.Title, pin.Description, pin.AllowComments, pin.ClickUrl, pin.ContentUrl,
		pin.LikesCount, pin.IsOwner, pin.IsLiked, pin.PinAuthor.UserId, pin.PinAuthor.Nickname, pin.PinAuthor.AvatarURL)
}

var PinPageResponses = []entity.PinPageResponse{
	{
		PinId:         entity.PinID(1),
		Title:         "title",
		Description:   "desc",
		ClickUrl:      "ccurl",
		IsOwner:       false,
		IsLiked:       false,
		AllowComments: true,
		LikesCount:    14,
		ContentUrl:    "imgjpg.com",
		PinAuthor: entity.PinAuthor{
			UserId:    entity.UserID(2),
			Nickname:  "NICK001",
			AvatarURL: "",
		},
	},
}

var pins = []entity.Pin{
	{
		PinId:         entity.PinID(1),
		AuthorId:      entity.UserID(2),
		ContentUrl:    "imgjpg.com",
		Title:         "title",
		Description:   "desc",
		AllowComments: true,
		ClickUrl:      "ccurl",
	},
}

func TestGetPin(t *testing.T) {
	type mockArgs struct {
		Ctx    context.Context
		PinId  entity.PinID
		UserId entity.UserID
		Slug   any
		Times  int
	}
	type mockReturn struct {
		Pin entity.PinPageResponse
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
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   1,
				Times:  1,
			},
			MockReturn: mockReturn{
				Pin: PinPageResponses[0],
				Err: errs.ErrorInfo{},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeGetPinResponse(PinPageResponses[0]),
				Code: 200,
			},
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   1,
				Times:  1,
			},
			MockReturn: mockReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: errs.ErrDBInternal},
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
		},
		{
			Name: "Uncorrect test 2",
			MockArgs: mockArgs{
				Ctx:    context.Background(),
				PinId:  entity.PinID(1),
				UserId: entity.UserID(0),
				Slug:   "abc",
				Times:  0,
			},
			MockReturn: mockReturn{
				Pin: entity.PinPageResponse{},
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal, GeneralErr: errs.ErrDBInternal},
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
		r := httptest.NewRequest(http.MethodGet, "/api/v1/pins/", nil)
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.MockArgs.Slug))
		w := httptest.NewRecorder()
		serviceMock.EXPECT().GetPinById(curTest.MockArgs.Ctx, curTest.MockArgs.PinId, curTest.MockArgs.UserId).
			Return(curTest.MockReturn.Pin, curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.GetPin(w, r)
		assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
		assert.Equal(t, curTest.ExpectedResponse.Body, w.Body.String())
	}
}

func TestUpdatePin(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		Pin   entity.Pin
		Slug  any
		Times int
	}
	type mockReturn struct {
		Pin entity.PinPageResponse
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
		Request          map[string]string
		ExpectedResponse expectedResponse
	}
	ctx := context.WithValue(context.Background(), "user_id", entity.UserID(0))
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Ctx: ctx,
				Pin: entity.Pin{
					Title: "title",
				},
				Slug:  1,
				Times: 1,
			},
			MockReturn: mockReturn{
				Pin: PinPageResponses[0],
				Err: errs.ErrorInfo{},
			},
			Request: map[string]string{
				"title": pins[0].Title,
			},
			ExpectedResponse: expectedResponse{
				Body: MakeGetPinResponse(PinPageResponses[0]),
				Code: 200,
			},
		},
		{
			Name: "Uncorrect test 1",
			MockArgs: mockArgs{
				Ctx: ctx,
				Pin: entity.Pin{
					Title: "title",
				},
				Slug:  1,
				Times: 1,
			},
			MockReturn: mockReturn{
				Pin: PinPageResponses[0],
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			Request: map[string]string{
				"title": pins[0].Title,
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
		},
		{
			Name: "Uncorrect test 2",
			MockArgs: mockArgs{
				Ctx: ctx,
				Pin: entity.Pin{
					Title: "title",
				},
				Slug:  "abc",
				Times: 0,
			},
			MockReturn: mockReturn{
				Pin: PinPageResponses[0],
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			Request: map[string]string{
				"title": pins[0].Title,
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
		reqBytes, err := json.Marshal(curTest.MockArgs.Pin)
		if err != nil {
			t.Errorf("error marshaling request body: %v", err)
		}
		r := httptest.NewRequest(http.MethodPost, "/api/v1/pins/", bytes.NewBuffer(reqBytes))
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.MockArgs.Slug))
		r = r.WithContext(curTest.MockArgs.Ctx)
		w := httptest.NewRecorder()
		serviceMock.EXPECT().UpdatePin(curTest.MockArgs.Ctx, curTest.MockArgs.Pin).
			Return(curTest.MockReturn.Pin, curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.UpdatePin(w, r)
		assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
		assert.Equal(t, curTest.ExpectedResponse.Body, w.Body.String())
	}
}

func TestDeletePin(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		Pin   entity.Pin
		Slug  any
		Times int
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
	ctx := context.WithValue(context.Background(), "user_id", entity.UserID(1))
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Ctx: ctx,
				Pin: entity.Pin{
					PinId:    entity.PinID(1),
					AuthorId: entity.UserID(1),
				},
				Slug:  1,
				Times: 1,
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
				Ctx: ctx,
				Pin: entity.Pin{
					PinId:    entity.PinID(1),
					AuthorId: entity.UserID(1),
				},
				Slug:  1,
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
				Ctx: ctx,
				Pin: entity.Pin{
					PinId:    entity.PinID(1),
					AuthorId: entity.UserID(1),
				},
				Slug:  "abc",
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
		r := httptest.NewRequest(http.MethodDelete, "/api/v1/pins/", nil)
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.MockArgs.Slug))
		r = r.WithContext(curTest.MockArgs.Ctx)
		w := httptest.NewRecorder()
		serviceMock.EXPECT().DeletePin(curTest.MockArgs.Ctx, curTest.MockArgs.Pin).
			Return(curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.DeletePin(w, r)
		assert.Equal(t, curTest.ExpectedResponse.Code, w.Code)
		assert.Equal(t, curTest.ExpectedResponse.Body, w.Body.String())
	}
}
