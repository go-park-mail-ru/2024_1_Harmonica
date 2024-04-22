package handler

import (
	"encoding/json"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"net/http"

	"go.uber.org/zap"
)

func WriteDefaultResponse(w http.ResponseWriter, logger *zap.Logger, object any) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(object)
	_, err := w.Write(response)
	if err != nil {
		logger.Error(
			errs.ErrServerInternal.Error(),
			zap.Int("local_error_code", errs.ErrorCodes[errs.ErrServerInternal].LocalCode),
			zap.String("general_error", err.Error()),
		)
	}
}

func WriteUserResponse(w http.ResponseWriter, logger *zap.Logger, user entity.User) {
	w.Header().Set("Content-Type", "application/json")
	userResponse := entity.UserResponse{
		UserId:    user.UserID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}
	response, _ := json.Marshal(userResponse)
	_, err := w.Write(response)
	if err != nil {
		logger.Error(
			errs.ErrServerInternal.Error(),
			zap.Int("local_error_code", errs.ErrorCodes[errs.ErrServerInternal].LocalCode),
			zap.String("general_error", err.Error()),
		)
	}
}
