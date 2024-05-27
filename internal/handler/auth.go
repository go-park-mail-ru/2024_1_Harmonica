package handler

import (
	"errors"
	"fmt"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	auth "harmonica/internal/microservices/auth/proto"
	"io"
	"net/http"
	"time"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var (
	emptyUser      = entity.User{}
	emptyErrorInfo = errs.ErrorInfo{}
)

// Login
//
//	@Summary		Login user
//	@Description	Login user by request.body json
//	@Tags			Authorization
//	@Produce		json
//	@Accept			json
//	@Header			200		{string}	Set-Cookie	"session-token"
//	@Param			Cookie	header		string		false	"session-token"	default(session-token=)
//	@Success		200		{object}	entity.User
//	@Failure		400		{object}	errs.ErrorResponse	"Possible code responses: 3, 4, 5."
//	@Failure		401		{object}	errs.ErrorResponse	"Possible code responses: 7, 8."
//	@Failure		403		{object}	errs.ErrorResponse	"Possible code responses: 1."
//	@Failure		500		{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/login [post]
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody,
		})
		return
	}

	var user entity.User
	err = easyjson.Unmarshal(bodyBytes, &user)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadingRequestBody})
		return
	}
	res, err := h.AuthService.Login(metadata.NewOutgoingContext(r.Context(),
		metadata.Pairs("request_id", requestId)),
		&auth.LoginUserRequest{
			UserId:   int64(user.UserID),
			Email:    user.Email,
			Nickname: user.Nickname,
			Password: user.Password,
		})
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.ErrUnauthorized, GeneralErr: err})
		return
	}
	if !res.Valid || res.LocalError != 0 {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.GetLocalErrorByCode[res.LocalError]})
		return
	}

	time, err := time.Parse(time.RFC3339Nano, res.ExpiresAt)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr: errs.ErrCantParseTime,
		})
		return
	}
	SetSessionTokenCookie(w, res.NewSessionToken, time)
	WriteUserResponse(w, h.logger, entity.User{
		UserID:    entity.UserID(res.UserId),
		Email:     res.Email,
		Nickname:  res.Nickname,
		Password:  res.Password,
		AvatarURL: res.AvatarURL,
	})
}

// Logout
//
//	@Summary		Logout user
//	@Description	Logout user by their session cookie
//	@Tags			Authorization
//	@Param			Cookie	header		string		false	"session-token"	default(session-token=)
//	@Header			200		{string}	Set-Cookie	"session-token"
//	@Success		200
//	@Failure		400	{object}	errs.ErrorResponse	"Possible code responses: 3."
//	@Router			/logout [get]
func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	c, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrReadCookie,
		})
		return
	}
	sessionToken := c.Value
	h.AuthService.Logout(r.Context(), &auth.LogoutRequest{SessionToken: sessionToken})
	SetSessionTokenCookie(w, sessionToken, time.Now())
	w.WriteHeader(http.StatusNoContent)
}

// Check if user is authorized
//
//	@Summary		Get auth status
//	@Description	Get user by request cookie
//	@Tags			Authorization
//	@Param			Cookie	header	string	false	"session-token"	default(session-token=)
//	@Produce		json
//	@Success		200	{object}	entity.UserResponse
//	@Failure		400	{object}	errs.ErrorResponse	"Possible code responses: 3."
//	@Failure		401	{object}	errs.ErrorResponse	"Possible code responses: 2."
//	@Failure		500	{object}	errs.ErrorResponse	"Possible code responses: 11."
//	@Router			/is_auth [get]
func (h *APIHandler) IsAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, ok1 := ctx.Value("request_id").(string)
	userIdFromSession, ok2 := ctx.Value("user_id").(entity.UserID)
	if !ok1 || !ok2 {
		WriteErrorResponse(w, h.logger, requestId, MakeErrorInfo(nil, errs.ErrUnauthorized))
		return
	}

	res, err := h.AuthService.IsAuth(metadata.NewOutgoingContext(r.Context(),
		metadata.Pairs("user_id", fmt.Sprintf(`%d`, userIdFromSession), "request_id", requestId)), &auth.Empty{})
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{
			LocalErr:   errs.ErrGRPCWentWrong,
			GeneralErr: err,
		})
		return
	}
	if !res.Valid {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.GetLocalErrorByCode[res.LocalError]})
		return
	}
	WriteUserResponse(w, h.logger, entity.User{
		UserID:    entity.UserID(res.User.UserId),
		Email:     res.User.Email,
		Nickname:  res.User.Nickname,
		AvatarURL: res.User.AvatarURL,
	})
}

func WriteUserResponse(w http.ResponseWriter, logger *zap.Logger, user entity.User) {
	w.Header().Set("Content-Type", "application/json")
	userResponse := entity.UserResponse{
		UserId:    user.UserID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}
	response, _ := easyjson.Marshal(userResponse)
	_, err := w.Write(response)
	if err != nil {
		logger.Error(
			errs.ErrServerInternal.Error(),
			zap.Int("local_error_code", errs.ErrorCodes[errs.ErrServerInternal].LocalCode),
			zap.String("general_error", err.Error()),
		)
	}
}

func SetSessionTokenCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
}
