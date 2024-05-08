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

var boards = []entity.Board{
	{
		BoardID:        1,
		Title:          "Board #1",
		Description:    "Public board number 1",
		CoverURL:       "",
		VisibilityType: "public",
	},
	{
		BoardID:        2,
		Title:          "Board #2",
		Description:    "Private board number 2",
		CoverURL:       "alala",
		VisibilityType: "private",
	},
	{
		BoardID:        3,
		Title:          "",
		Description:    "Wrong board (no title)",
		CoverURL:       "",
		VisibilityType: "public",
	},
}

var boardAuthor = entity.BoardAuthor{
	UserId:   1,
	Nickname: "MaryPoppins",
}

var pinAuthor = entity.PinAuthor{
	UserId:   1,
	Nickname: "MaryPoppins",
}

var boardPinResponse = entity.BoardPinResponse{
	PinId:      1,
	ContentUrl: "pin.jpg",
	PinAuthor:  pinAuthor,
}

var fullBoards = []entity.FullBoard{
	{
		Board:        boards[0],
		BoardAuthors: []entity.BoardAuthor{boardAuthor},
		Pins:         []entity.BoardPinResponse{boardPinResponse},
	},
}

func MakeDefaultResponse(t *testing.T, response any) string {
	jsonBodyBytes, err := json.Marshal(response)
	if err != nil {
		t.Errorf("error marshalling response: %v", err)
	}
	jsonBody := string(jsonBodyBytes)
	return jsonBody
}

func TestCreateBoard(t *testing.T) {
	type mockArgs struct {
		Board  entity.Board
		UserId entity.UserID
		Times  int
	}
	type mockReturn struct {
		FullBoard entity.FullBoard
		Err       errs.ErrorInfo
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
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				Board:  boards[0],
				UserId: users[0].UserID,
				Times:  1,
			},
			MockReturn: mockReturn{
				FullBoard: fullBoards[0],
			},
			Request: []byte(fmt.Sprintf(`{
					"title": "%s",
					"description": "%s",
					"cover_url": "%s",
					"visibility_type": "%s"
				}`, boards[0].Title, boards[0].Description, boards[0].CoverURL, boards[0].VisibilityType)),
			ExpectedResponse: expectedResponse{
				Body: MakeDefaultResponse(t, fullBoards[0]),
				Code: 200,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name:       "Incorrect test 1",
			MockArgs:   mockArgs{},
			MockReturn: mockReturn{},
			Request: []byte(fmt.Sprintf(`{
					"title": "%s",
					"description": "%s",
					"cover_url": "%s",
					"visibility_type": "%s"
				}`, boards[2].Title, boards[2].Description, boards[2].CoverURL, boards[2].VisibilityType)),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidInputFormat),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name:     "Incorrect test 2",
			MockArgs: mockArgs{Times: 1},
			MockReturn: mockReturn{
				FullBoard: entity.FullBoard{},
				Err:       errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			Request: []byte(fmt.Sprintf(`{
					"title": "%s",
					"description": "%s",
					"cover_url": "%s",
					"visibility_type": "%s"
				}`, boards[0].Title, boards[0].Description, boards[0].CoverURL, boards[0].VisibilityType)),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name:       "Incorrect test 2",
			MockArgs:   mockArgs{},
			MockReturn: mockReturn{},
			Request:    []byte(`"title": alala"})`),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrReadingRequestBody),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/boards", bytes.NewBuffer(curTest.Request))
		w := httptest.NewRecorder()
		ctx := context.WithValue(curTest.Ctx, "user_id", curTest.MockArgs.UserId)
		r = r.WithContext(ctx)
		serviceMock.EXPECT().CreateBoard(ctx, gomock.Any(), curTest.MockArgs.UserId).
			Return(curTest.MockReturn.FullBoard, curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.CreateBoard(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}

func TestGetBoard(t *testing.T) {
	type mockArgs struct {
		BoardId entity.BoardID
		UserId  entity.UserID
	}
	type mockReturn struct {
		FullBoard entity.FullBoard
		Err       errs.ErrorInfo
	}
	type expectedResponse struct {
		Body string
		Code int
	}
	type test struct {
		Name             string
		MockArgs         mockArgs
		MockReturn       mockReturn
		Slug             any
		Query            string
		ExpectedResponse expectedResponse
		Ctx              context.Context
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				UserId:  users[0].UserID,
			},
			Slug: boards[0].BoardID,
			MockReturn: mockReturn{
				FullBoard: fullBoards[0],
			},
			ExpectedResponse: expectedResponse{
				Body: MakeDefaultResponse(t, fullBoards[0]),
				Code: 200,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				UserId:  users[0].UserID,
			},
			MockReturn: mockReturn{
				FullBoard: entity.FullBoard{},
				Err:       errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			Slug: boards[0].BoardID,
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 2",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				UserId:  users[0].UserID,
			},
			MockReturn: mockReturn{},
			Slug:       "abc",
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name:       "Incorrect test 3",
			MockArgs:   mockArgs{},
			MockReturn: mockReturn{},
			Slug:       1,
			Query:      "lala",
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrReadingRequestBody),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/boards/{board_id}", nil)
		w := httptest.NewRecorder()
		r.SetPathValue("board_id", fmt.Sprintf(`%d`, curTest.Slug))
		if len(curTest.Query) != 0 {
			q := r.URL.Query()
			q.Set("page", curTest.Query)
			r.URL.RawQuery = q.Encode()
		}
		ctx := context.WithValue(curTest.Ctx, "user_id", curTest.MockArgs.UserId)
		r = r.WithContext(ctx)
		serviceMock.EXPECT().GetBoardById(ctx, curTest.MockArgs.BoardId, curTest.MockArgs.UserId,
			gomock.Any(), gomock.Any()).Return(curTest.MockReturn.FullBoard, curTest.MockReturn.Err).MaxTimes(1)
		h.GetBoard(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}

func TestAddPinToBoard(t *testing.T) {
	type mockArgs struct {
		BoardId entity.BoardID
		PinId   entity.PinID
		UserId  entity.UserID
		Times   int
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
		SlugPinId        any
		SlugBoardId      any
		ExpectedResponse expectedResponse
		Ctx              context.Context
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
				Times:   1,
			},
			MockReturn:  mockReturn{},
			SlugPinId:   entity.PinID(1),
			SlugBoardId: boards[0].BoardID,
			ExpectedResponse: expectedResponse{
				Body: "null",
				Code: 200,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
			},
			MockReturn: mockReturn{},
			SlugPinId:  "alala",
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 2",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
			},
			MockReturn:  mockReturn{},
			SlugBoardId: "lalala",
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 3",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
				Times:   1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			SlugPinId:   entity.PinID(1),
			SlugBoardId: boards[0].BoardID,
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/boards/{board_id}/pins/{pin_id}", nil)
		w := httptest.NewRecorder()
		r.SetPathValue("board_id", fmt.Sprintf(`%d`, curTest.SlugPinId))
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.SlugBoardId))
		ctx := context.WithValue(curTest.Ctx, "user_id", curTest.MockArgs.UserId)
		r = r.WithContext(ctx)
		serviceMock.EXPECT().AddPinToBoard(ctx, curTest.MockArgs.BoardId, curTest.MockArgs.PinId,
			curTest.MockArgs.UserId).Return(curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.AddPinToBoard(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}

func TestDeletePinFromBoard(t *testing.T) {
	type mockArgs struct {
		BoardId entity.BoardID
		PinId   entity.PinID
		UserId  entity.UserID
		Times   int
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
		SlugPinId        any
		SlugBoardId      any
		ExpectedResponse expectedResponse
		Ctx              context.Context
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
				Times:   1,
			},
			MockReturn:  mockReturn{},
			SlugPinId:   entity.PinID(1),
			SlugBoardId: boards[0].BoardID,
			ExpectedResponse: expectedResponse{
				Body: "null",
				Code: 200,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
			},
			MockReturn: mockReturn{},
			SlugPinId:  "alala",
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 2",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
			},
			MockReturn:  mockReturn{},
			SlugBoardId: "lalala",
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 3",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				PinId:   entity.PinID(1),
				UserId:  users[0].UserID,
				Times:   1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrElementNotExist},
			},
			SlugPinId:   entity.PinID(1),
			SlugBoardId: boards[0].BoardID,
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrElementNotExist),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/boards/{board_id}/pins/{pin_id}", nil)
		w := httptest.NewRecorder()
		r.SetPathValue("board_id", fmt.Sprintf(`%d`, curTest.SlugBoardId))
		r.SetPathValue("pin_id", fmt.Sprintf(`%d`, curTest.SlugPinId))
		ctx := context.WithValue(curTest.Ctx, "user_id", curTest.MockArgs.UserId)
		r = r.WithContext(ctx)
		serviceMock.EXPECT().DeletePinFromBoard(ctx, curTest.MockArgs.BoardId, curTest.MockArgs.PinId,
			curTest.MockArgs.UserId).Return(curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.DeletePinFromBoard(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}

func TestDeleteBoard(t *testing.T) {
	type mockArgs struct {
		BoardId entity.BoardID
		UserId  entity.UserID
		Times   int
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
		Slug             any
		ExpectedResponse expectedResponse
		Ctx              context.Context
	}
	tests := []test{
		{
			Name: "Correct test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				UserId:  users[0].UserID,
				Times:   1,
			},
			Slug:       boards[0].BoardID,
			MockReturn: mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: "null",
				Code: 200,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 1",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				UserId:  users[0].UserID,
				Times:   1,
			},
			MockReturn: mockReturn{
				Err: errs.ErrorInfo{LocalErr: errs.ErrDBInternal},
			},
			Slug: boards[0].BoardID,
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrDBInternal),
				Code: 500,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
		{
			Name: "Incorrect test 2",
			MockArgs: mockArgs{
				BoardId: boards[0].BoardID,
				UserId:  users[0].UserID,
			},
			MockReturn: mockReturn{},
			Slug:       "abc",
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidSlug),
				Code: 400,
			},
			Ctx: context.WithValue(context.Background(), "request_id", "req_id"),
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L(), nil, nil, nil, nil)
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodDelete, "/api/v1/boards/", nil)
		w := httptest.NewRecorder()
		r.SetPathValue("board_id", fmt.Sprintf(`%d`, curTest.Slug))
		ctx := context.WithValue(curTest.Ctx, "user_id", curTest.MockArgs.UserId)
		r = r.WithContext(ctx)
		serviceMock.EXPECT().DeleteBoard(ctx, curTest.MockArgs.BoardId, curTest.MockArgs.UserId).
			Return(curTest.MockReturn.Err).Times(curTest.MockArgs.Times)
		h.DeleteBoard(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}
