package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/handler"
	mock_service "harmonica/mocks/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

var users = []entity.User{
	{
		UserID:   1,
		Email:    "mary@email.ru",
		Nickname: "MaryPoppins",
		Password: "MaryPoppins25",
	},
	{
		UserID:   2,
		Email:    "TesT.her_e@sth.ru",
		Nickname: "Michael",
		Password: "1918Michael",
	},
	{
		UserID:   3,
		Email:    "Crazy%e.{m}ail~@$trange.com",
		Nickname: "crazy_user",
		Password: "crazyUser24ksHokssn27awb",
	},
	{
		UserID:   4,
		Email:    "something_wrong",
		Nickname: "crazy_user",
		Password: "password",
	},
}

//type testInfo struct {
//	Name string
//	Code int
//}
//loginTestsNames := []{
//"Correct test 1", "Correct test 2", "Correct test 3",
//}
//
//loginTestsCodes := []int {
//200, 200, 200
//}

func MakeUserResponseBody(user entity.User) string {
	return fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s","avatar_url":"%s"}`,
		user.UserID, user.Email, user.Nickname, user.AvatarURL)
}

func MakeErrorResponse(err error) string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`,
		errs.ErrorCodes[err].LocalCode, err.Error())
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func MakeErrorListResponse(errsList ...error) string {
	errors := make([]Error, len(errsList))
	for i, err := range errsList {
		errors[i] = Error{
			Code:    errs.ErrorCodes[err].LocalCode,
			Message: err.Error(),
		}
	}
	response := ErrorResponse{Errors: errors}
	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

func TestLogin(t *testing.T) {
	type mockArgs struct {
		Ctx   context.Context
		Email string
	}
	type mockReturn struct {
		User entity.User
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
		Request          []byte
		ExpectedResponse expectedResponse
	}
	tests := []test{
		{
			Name: "Correct request 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: users[0].Email,
			},
			MockReturn: mockReturn{
				User: users[0],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[0].Email, users[0].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeUserResponseBody(users[0]),
				Code: 200,
			},
		},
		{
			Name: "Correct request 2",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: users[1].Email,
			},
			MockReturn: mockReturn{
				User: users[1],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[1].Email, users[1].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeUserResponseBody(users[1]),
				Code: 200,
			},
		},
		{
			Name: "Correct request 3",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: users[2].Email,
			},
			MockReturn: mockReturn{
				User: users[2],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[2].Email, users[2].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeUserResponseBody(users[2]),
				Code: 200,
			},
		},
		{
			Name: "Incorrect request 1",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: users[3].Email,
			},
			MockReturn: mockReturn{
				User: users[3],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, users[3].Email, users[3].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidInputFormat),
				Code: 400,
			},
		},
		{
			Name: "Incorrect request 2",
			MockArgs: mockArgs{
				Ctx:   context.Background(),
				Email: users[0].Email,
			},
			MockReturn: mockReturn{
				User: users[0],
			},
			Request: []byte(`"alala`),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrReadingRequestBody),
				Code: 400,
			},
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L())
	for _, curTest := range tests {
		curHashedPassword, err := bcrypt.GenerateFromPassword([]byte(curTest.MockReturn.User.Password), bcrypt.DefaultCost)
		if err != nil {
			t.Errorf("error hashng passwords: %v", err)
			return
		}
		curTest.MockReturn.User.Password = string(curHashedPassword)
		r := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(curTest.Request))
		w := httptest.NewRecorder()
		serviceMock.EXPECT().GetUserByEmail(curTest.MockArgs.Ctx, curTest.MockArgs.Email).
			Return(curTest.MockReturn.User, curTest.MockReturn.Err).MaxTimes(1)
		h.Login(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}

func TestLogout(t *testing.T) {
	type expectedResponse struct {
		Body string
		Code int
	}
	type test struct {
		Name             string
		ExpectedResponse expectedResponse
		Cookie           *http.Cookie
	}
	tests := []test{
		{
			Name: "Correct request 1",
			ExpectedResponse: expectedResponse{
				Code: 200,
			},
		},
		{
			Name: "Correct request 1",
			ExpectedResponse: expectedResponse{
				Code: 200,
			},
			Cookie: &http.Cookie{
				Name:  "session_token",
				Value: "token",
			},
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L())
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
		if curTest.Cookie != nil {
			r.AddCookie(curTest.Cookie)
		}
		w := httptest.NewRecorder()
		h.Logout(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}

func TestRegister(t *testing.T) {
	type mockArgs struct {
		Ctx  context.Context
		User entity.User
	}
	type mockReturn struct {
		User         entity.User
		RegisterErrs []errs.ErrorInfo
		GetUserErr   errs.ErrorInfo
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
	}
	tests := []test{
		{
			Name: "Correct request 1",
			MockArgs: mockArgs{
				Ctx:  context.Background(),
				User: users[0],
			},
			MockReturn: mockReturn{
				User: users[0],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","nickname":"%s","password":"%s"}`,
				users[0].Email, users[0].Nickname, users[0].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeUserResponseBody(users[0]),
				Code: 200,
			},
		},
		{
			Name: "Correct request 2",
			MockArgs: mockArgs{
				Ctx:  context.Background(),
				User: users[1],
			},
			MockReturn: mockReturn{
				User: users[1],
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","nickname":"%s","password":"%s"}`,
				users[1].Email, users[1].Nickname, users[1].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeUserResponseBody(users[1]),
				Code: 200,
			},
		},
		{
			Name: "Incorrect request 1",
			MockArgs: mockArgs{
				Ctx:  context.Background(),
				User: users[2],
			},
			MockReturn: mockReturn{
				RegisterErrs: []errs.ErrorInfo{
					{LocalErr: errs.ErrDBUniqueEmail},
					{LocalErr: errs.ErrDBUniqueEmail},
				},
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","nickname":"%s","password":"%s"}`,
				users[2].Email, users[2].Nickname, users[2].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorListResponse(errs.ErrDBUniqueEmail, errs.ErrDBUniqueEmail),
				Code: 500,
			},
		},
		{
			Name: "Incorrect request 2",
			MockArgs: mockArgs{
				Ctx:  context.Background(),
				User: users[0],
			},
			MockReturn: mockReturn{
				GetUserErr: errs.ErrorInfo{LocalErr: errs.ErrUserNotExist},
			},
			Request: []byte(fmt.Sprintf(`{"email":"%s","nickname":"%s","password":"%s"}`,
				users[0].Email, users[0].Nickname, users[0].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorListResponse(errs.ErrUserNotExist),
				Code: 404,
			},
		},
		{
			Name:       "Incorrect request 3",
			MockArgs:   mockArgs{},
			MockReturn: mockReturn{},
			Request: []byte(fmt.Sprintf(`{"email":"%s","nickname":"%s","password":"%s"}`,
				users[3].Email, users[3].Nickname, users[3].Password)),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorListResponse(errs.ErrInvalidInputFormat),
				Code: 400,
			},
		},
		{
			Name:       "Incorrect request 4",
			MockArgs:   mockArgs{},
			MockReturn: mockReturn{},
			Request:    []byte(`{"blabla"")`),
			ExpectedResponse: expectedResponse{
				Body: MakeErrorListResponse(errs.ErrReadingRequestBody),
				Code: 400,
			},
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L())
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(curTest.Request))
		w := httptest.NewRecorder()
		serviceMock.EXPECT().RegisterUser(curTest.MockArgs.Ctx, gomock.Any()).
			Return(curTest.MockReturn.RegisterErrs).MaxTimes(1)
		serviceMock.EXPECT().GetUserByEmail(curTest.MockArgs.Ctx, curTest.MockArgs.User.Email).
			Return(curTest.MockReturn.User, curTest.MockReturn.GetUserErr).MaxTimes(1)
		h.Register(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}

func TestIsAuth(t *testing.T) {
	type mockArgs struct {
		Ctx  context.Context
		User entity.User
	}
	type mockReturn struct {
		User entity.User
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
		RequestCtx       context.Context
		ExpectedResponse expectedResponse
	}
	tests := []test{
		{
			Name: "Correct request 1",
			MockArgs: mockArgs{
				//Ctx:  context.WithValue(context.Background(), "user_id", users[0].UserID),
				Ctx:  context.Background(),
				User: users[0],
			},
			MockReturn: mockReturn{
				User: users[0],
			},
			//RequestCtx: context.WithValue(context.Background(), "user_id", users[0].UserID),
			ExpectedResponse: expectedResponse{
				Body: MakeUserResponseBody(users[0]),
				Code: 200,
			},
		},
		{
			Name: "Incorrect request 1",
			MockArgs: mockArgs{
				Ctx: context.Background(),
			},
			MockReturn: mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrUnauthorized),
				Code: 401,
			},
		},
		{
			Name: "Incorrect request 2",
			MockArgs: mockArgs{
				Ctx:  context.Background(),
				User: users[0],
			},

			MockReturn: mockReturn{},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrUnauthorized),
				Code: 401,
			},
		},
	}
	ctrl := gomock.NewController(t)
	serviceMock := mock_service.NewMockIService(ctrl)
	h := handler.NewAPIHandler(serviceMock, zap.L())
	for _, curTest := range tests {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/is_auth", nil)
		w := httptest.NewRecorder()
		ctx := context.WithValue(context.Background(), "user_id", curTest.MockArgs.User.UserID)
		r = r.WithContext(ctx)
		serviceMock.EXPECT().GetUserById(ctx, curTest.MockArgs.User.UserID).
			Return(curTest.MockReturn.User, curTest.MockReturn.Err).MaxTimes(1)
		h.IsAuth(w, r)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}
