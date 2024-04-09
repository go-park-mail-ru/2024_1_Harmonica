package tests

//
//import (
//	"bytes"
//	"context"
//	"fmt"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"go.uber.org/zap"
//	"golang.org/x/crypto/bcrypt"
//	"harmonica/internal/entity"
//	"harmonica/internal/entity/errs"
//	"harmonica/internal/handler"
//	mock_service "harmonica/mocks/service"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestCreateBoard(t *testing.T) {
//	type mockArgs struct {
//		Ctx   context.Context
//		Board entity.Board
//	}
//	type mockReturn struct {
//		User entity.User
//		Err  errs.ErrorInfo
//	}
//	type request struct {
//	}
//	type expectedResponse struct {
//		Body string
//		Code int
//	}
//	type test struct {
//		Name             string
//		MockArgs         mockArgs
//		MockReturn       mockReturn
//		Request          []byte
//		ExpectedResponse expectedResponse
//	}
//	tests := []test{
//		{
//			Name: "Correct request 1",
//			MockArgs: mockArgs{
//				Ctx:   context.Background(),
//				Email: users[0].Email,
//			},
//			MockReturn: mockReturn{
//				User: users[0],
//			},
//			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[0].Email, users[0].Password)),
//			ExpectedResponse: expectedResponse{
//				Body: MakeUserResponseBody(users[0]),
//				Code: 200,
//			},
//		},
//		{
//			Name: "Correct request 2",
//			MockArgs: mockArgs{
//				Ctx:   context.Background(),
//				Email: users[1].Email,
//			},
//			MockReturn: mockReturn{
//				User: users[1],
//			},
//			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[1].Email, users[1].Password)),
//			ExpectedResponse: expectedResponse{
//				Body: MakeUserResponseBody(users[1]),
//				Code: 200,
//			},
//		},
//		{
//			Name: "Correct request 3",
//			MockArgs: mockArgs{
//				Ctx:   context.Background(),
//				Email: users[2].Email,
//			},
//			MockReturn: mockReturn{
//				User: users[2],
//			},
//			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[2].Email, users[2].Password)),
//			ExpectedResponse: expectedResponse{
//				Body: MakeUserResponseBody(users[2]),
//				Code: 200,
//			},
//		},
//		{
//			Name: "Incorrect request 1",
//			MockArgs: mockArgs{
//				Ctx:   context.Background(),
//				Email: users[3].Email,
//			},
//			MockReturn: mockReturn{
//				User: users[3],
//			},
//			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[3].Email, users[3].Password)),
//			ExpectedResponse: expectedResponse{
//				Body: MakeErrorResponse(errs.ErrInvalidInputFormat),
//				Code: 400,
//			},
//		},
//		{
//			Name: "Incorrect request 2",
//			MockArgs: mockArgs{
//				Ctx:   context.Background(),
//				Email: users[0].Email,
//			},
//			MockReturn: mockReturn{
//				User: users[0],
//			},
//			Request: []byte(`"alala`),
//			ExpectedResponse: expectedResponse{
//				Body: MakeErrorResponse(errs.ErrReadingRequestBody),
//				Code: 400,
//			},
//		},
//	}
//	ctrl := gomock.NewController(t)
//	serviceMock := mock_service.NewMockIService(ctrl)
//	h := handler.NewAPIHandler(serviceMock, zap.L())
//	for _, curTest := range tests {
//		curHashedPassword, err := bcrypt.GenerateFromPassword([]byte(curTest.MockReturn.User.Password), bcrypt.DefaultCost)
//		if err != nil {
//			t.Errorf("error hashng passwords: %v", err)
//			return
//		}
//		curTest.MockReturn.User.Password = string(curHashedPassword)
//		r := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(curTest.Request))
//		w := httptest.NewRecorder()
//		serviceMock.EXPECT().GetUserByEmail(curTest.MockArgs.Ctx, curTest.MockArgs.Email).
//			Return(curTest.MockReturn.User, curTest.MockReturn.Err).MaxTimes(1)
//		h.Login(w, r)
//		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
//		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
//	}
//}
