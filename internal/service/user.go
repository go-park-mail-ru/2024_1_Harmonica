package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

var emptyUser = entity.User{}

func (s *RepositoryService) GetUserByEmail(ctx context.Context, email string) (entity.User, errs.ErrorInfo) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, emptyErrorInfo
}

func (s *RepositoryService) GetUserByNickname(ctx context.Context, nickname string) (entity.User, errs.ErrorInfo) {
	user, err := s.repo.GetUserByNickname(ctx, nickname)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, emptyErrorInfo
}

func (s *RepositoryService) GetUserById(ctx context.Context, id entity.UserID) (entity.User, errs.ErrorInfo) {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, emptyErrorInfo
}

func (s *RepositoryService) RegisterUser(ctx context.Context, user entity.User) []errs.ErrorInfo {
	user.Sanitize()
	var errsList []errs.ErrorInfo

	// Checking for unique fields
	checkUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		errsList = append(errsList, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		})
		return errsList
	}
	if checkUser != emptyUser {
		errsList = append(errsList, errs.ErrorInfo{
			GeneralErr: nil,
			LocalErr:   errs.ErrDBUniqueEmail,
		})
	}

	checkUser, err = s.repo.GetUserByNickname(ctx, user.Nickname)
	if err != nil {
		errsList = append(errsList, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		})
		return errsList
	}
	if checkUser != emptyUser {
		errsList = append(errsList, errs.ErrorInfo{
			GeneralErr: nil,
			LocalErr:   errs.ErrDBUniqueNickname},
		)
	}

	if len(errsList) > 0 {
		return errsList
	}

	err = s.repo.RegisterUser(ctx, user)
	if err != nil {
		errsList = append(errsList, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		})
	}

	return errsList
}

func (s *RepositoryService) UpdateUser(ctx context.Context, user entity.User) (entity.User, errs.ErrorInfo) {
	user.Sanitize()
	if user.Nickname != "" {
		checkUser, err := s.repo.GetUserByNickname(ctx, user.Nickname)
		if err != nil {
			return emptyUser, errs.ErrorInfo{
				GeneralErr: err,
				LocalErr:   errs.ErrDBInternal,
			}
		}
		if checkUser != emptyUser && checkUser.UserID != user.UserID {
			return emptyUser, errs.ErrorInfo{
				GeneralErr: nil,
				LocalErr:   errs.ErrDBUniqueNickname,
			}
		}
	}

	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	updatedUser, err := s.repo.GetUserById(ctx, user.UserID)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	return updatedUser, emptyErrorInfo
}

func (s *RepositoryService) GetUserProfileByNickname(ctx context.Context, nickname string,
	userId entity.UserID) (entity.UserProfileResponse, errs.ErrorInfo) {
	user, err := s.repo.GetUserByNickname(ctx, nickname)
	if err != nil {
		return entity.UserProfileResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	userProfile, err := s.repo.GetSubscriptionsInfo(ctx, user.UserID, userId)
	if err != nil {
		return entity.UserProfileResponse{}, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	userProfile.User = MakeUserResponse(user)
	userProfile.IsOwner = userProfile.User.UserId == userId
	return userProfile, emptyErrorInfo
}

func MakeUserResponse(user entity.User) entity.UserResponse {
	userResponse := entity.UserResponse{
		UserId:    user.UserID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		AvatarDX:  user.AvatarDX,
		AvatarDY:  user.AvatarDY,
	}
	return userResponse
}
