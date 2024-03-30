package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

var emptyUser = entity.User{}

func (r *RepositoryService) GetUserByEmail(ctx context.Context, email string) (entity.User, errs.ErrorInfo) {
	user, err := r.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, emptyErrorInfo
}

func (r *RepositoryService) GetUserByNickname(ctx context.Context, nickname string) (entity.User, errs.ErrorInfo) {
	user, err := r.repo.GetUserByNickname(ctx, nickname)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, emptyErrorInfo
}

func (r *RepositoryService) GetUserById(ctx context.Context, id entity.UserID) (entity.User, errs.ErrorInfo) {
	user, err := r.repo.GetUserById(ctx, id)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}
	return user, emptyErrorInfo
}

func (r *RepositoryService) RegisterUser(ctx context.Context, user entity.User) []errs.ErrorInfo {
	var errsList []errs.ErrorInfo

	// Checking for unique fields
	checkUser, err := r.repo.GetUserByEmail(ctx, user.Email)
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

	checkUser, err = r.repo.GetUserByNickname(ctx, user.Nickname)
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

	err = r.repo.RegisterUser(ctx, user)
	if err != nil {
		errsList = append(errsList, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		})
	}

	return errsList
}

func (r *RepositoryService) UpdateUser(ctx context.Context, user entity.User) (entity.User, errs.ErrorInfo) {
	if user.Nickname != "" {
		checkUser, err := r.repo.GetUserByNickname(ctx, user.Nickname)
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

	err := r.repo.UpdateUser(ctx, user)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	updatedUser, err := r.repo.GetUserById(ctx, user.UserID)
	if err != nil {
		return emptyUser, errs.ErrorInfo{
			GeneralErr: err,
			LocalErr:   errs.ErrDBInternal,
		}
	}

	return updatedUser, emptyErrorInfo
}
