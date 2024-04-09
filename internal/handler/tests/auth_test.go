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
		Email:    "something",
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

func MakeLoginResponseBody(user entity.User) string {
	return fmt.Sprintf(`{"user_id":%d,"email":"%s","nickname":"%s","avatar_url":"%s"}`,
		user.UserID, user.Email, user.Nickname, user.AvatarURL)
}

func MakeErrorResponse(err error) string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`,
		errs.ErrorCodes[err].LocalCode, err.Error())
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
	//type funcArgs struct {
	//	Ctx   context.Context
	//	Email string
	//}
	//type expectedResponse struct {
	//	User entity.User
	//	Err  errs.ErrorInfo
	//}
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
			Request: map[string]string{
				"email":    users[0].Email,
				"password": users[0].Password,
			},
			ExpectedResponse: expectedResponse{
				Body: MakeLoginResponseBody(users[0]),
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
			Request: map[string]string{
				"email":    users[1].Email,
				"password": users[1].Password,
			},
			ExpectedResponse: expectedResponse{
				Body: MakeLoginResponseBody(users[1]),
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
			Request: map[string]string{
				"email":    users[2].Email,
				"password": users[2].Password,
			},
			ExpectedResponse: expectedResponse{
				Body: MakeLoginResponseBody(users[2]),
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
			Request: map[string]string{
				"email":    users[3].Email,
				"password": users[3].Password,
			},
			ExpectedResponse: expectedResponse{
				Body: MakeErrorResponse(errs.ErrInvalidInputFormat),
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
		reqBytes, err := json.Marshal(curTest.Request)
		if err != nil {
			t.Errorf("error marshaling request body: %v", err)
		}
		curTest.MockReturn.User.Password = string(curHashedPassword)
		r := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBytes))
		w := httptest.NewRecorder()
		serviceMock.EXPECT().GetUserByEmail(curTest.MockArgs.Ctx, curTest.MockArgs.Email).
			Return(curTest.MockReturn.User, curTest.MockReturn.Err).AnyTimes()
		//h := handler.NewAPIHandler(s, zap.L())
		h.Login(w, r)
		//user, err := h.Login(test.FuncArgs.Ctx, test.FuncArgs.Email)
		assert.Equal(t, w.Code, curTest.ExpectedResponse.Code)
		assert.Equal(t, w.Body.String(), curTest.ExpectedResponse.Body)
	}
}
