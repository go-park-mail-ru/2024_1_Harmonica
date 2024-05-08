package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	auth "harmonica/internal/microservices/auth/proto"
	"testing"
)

func TestAuthorizationServer_CheckSession(t *testing.T) {
	type SessionData struct {
		SessionID string
		UserID    int
		Expired   bool
	}

	testCases := []struct {
		name           string
		request        *auth.CheckSessionRequest
		sessionData    map[string]SessionData
		expectedOutput *auth.CheckSessionResponse
	}{
		{
			name: "Expired session",
			request: &auth.CheckSessionRequest{
				Session: "expired_session_id",
			},
			sessionData: map[string]SessionData{
				"expired_session_id": {SessionID: "expired_session_id", UserID: 456, Expired: true},
			},
			expectedOutput: &auth.CheckSessionResponse{Valid: false, LocalError: 2},
		},
		{
			name: "Non-existent session",
			request: &auth.CheckSessionRequest{
				Session: "nonexistent_session_id",
			},
			sessionData:    map[string]SessionData{},
			expectedOutput: &auth.CheckSessionResponse{Valid: false, LocalError: 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := AuthorizationServer{}

			// Populate session data
			for sessionID, data := range tc.sessionData {
				if data.Expired {
					// Simulate session expiration
					Sessions.Delete(sessionID)
				} else {
					// Simulate active session
					Sessions.Store(sessionID, Session{UserId: entity.UserID(data.UserID)})
				}
			}

			// Execute the function being tested
			response, err := server.CheckSession(context.Background(), tc.request)

			// Verify the output
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, response)
		})
	}
}

/* // не получается передать userId в контекст
func TestLogin(t *testing.T) {
	type Args struct {
		Email    string
		Password string
		Context  context.Context
	}
	type ExpectedReturn struct {
		Response *auth.LoginUserResponse
		Error    error
	}
	type ExpectedMockArgs struct {
		Email    string
		Password string
		Context  context.Context
	}
	type ExpectedMockReturn struct {
		User      entity.User
		ErrorInfo errs.ErrorInfo
	}

	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		mockArgs ExpectedMockArgs, mockReturn ExpectedMockReturn) {
		service.EXPECT().GetUserByEmail(ctx, mockArgs.Email).Return(mockReturn.User, mockReturn.ErrorInfo)
	}

	testTable := []struct {
		name               string
		args               Args
		expectedReturn     ExpectedReturn
		expectedMockArgs   ExpectedMockArgs
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test",
			args: Args{
				Email:    "example@example.com",
				Password: "password",
				Context:  context.Background(),
			},
			expectedReturn: ExpectedReturn{
				Response: &auth.LoginUserResponse{
					Valid:           true,
					UserId:          1,
					Email:           "example@example.com",
					Nickname:        "example",
					Password:        "hashedPassword",
					AvatarURL:       "example.com/avatar",
					RegisterAt:      "2024-05-09T12:00:00Z",
					NewSessionToken: "newSessionToken",
					ExpiresAt:       "2024-05-09T13:00:00Z",
				},
			},
			expectedMockArgs: ExpectedMockArgs{
				Email:    "example@example.com",
				Password: "password",
				Context:  context.WithValue(context.Background(), "user_id", 1),
			},
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{
					UserID:     1,
					Email:      "example@example.com",
					Nickname:   "example",
					Password:   "hashedPassword",
					AvatarURL:  "example.com/avatar",
					RegisterAt: time.Date(2024, 5, 9, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Invalid email",
			args: Args{
				Email:    "invalidEmail",
				Password: "password",
				Context:  context.Background(),
			},
			expectedReturn: ExpectedReturn{
				Response: &auth.LoginUserResponse{
					LocalError: 5,
					Valid:      false,
				},
			},
			expectedMockArgs: ExpectedMockArgs{
				Email:    "invalidEmail",
				Password: "password",
				Context:  context.Background(),
			},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name: "Invalid password",
			args: Args{
				Email:    "example@example.com",
				Password: "invalidPassword",
				Context:  context.Background(),
			},
			expectedReturn: ExpectedReturn{
				Response: &auth.LoginUserResponse{
					LocalError: 8,
					Valid:      false,
				},
			},
			expectedMockArgs: ExpectedMockArgs{
				Email:    "example@example.com",
				Password: "invalidPassword",
				Context:  context.Background(),
			},
			expectedMockReturn: ExpectedMockReturn{
				User: entity.User{
					UserID:     1,
					Email:      "example@example.com",
					Nickname:   "example",
					Password:   "hashedPassword",
					AvatarURL:  "example.com/avatar",
					RegisterAt: time.Date(2024, 5, 9, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mock_service.NewMockIService(ctrl)
			server := NewAuthorizationServerForTests(service)

			md := metadata.Pairs("user_id", "1")
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			mockBehaviour(service, ctx, testCase.expectedMockArgs, testCase.expectedMockReturn)
			response, err := server.Login(ctx, &auth.LoginUserRequest{
				Email:    testCase.args.Email,
				Password: testCase.args.Password,
			})
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}
*/
