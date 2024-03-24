package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
)

var emptyUser = entity.User{}

func (r *RepositoryService) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := r.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return emptyUser, errs.ErrDBInternal
	}
	return user, nil
}

func (r *RepositoryService) GetUserByNickname(ctx context.Context, nickname string) (entity.User, error) {
	user, err := r.repo.GetUserByNickname(ctx, nickname)
	if err != nil {
		return emptyUser, errs.ErrDBInternal
	}
	return user, nil
}

func (r *RepositoryService) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	user, err := r.repo.GetUserById(ctx, id)
	if err != nil {
		return emptyUser, errs.ErrDBInternal
	}
	return user, nil
}

func (r *RepositoryService) RegisterUser(ctx context.Context, user entity.User) []error {
	var errsList []error

	// Checking for unique fields
	checkUser, err := r.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		errsList = append(errsList, errs.ErrDBInternal)
		return errsList
	}
	if checkUser != emptyUser {
		errsList = append(errsList, errs.ErrDBUniqueEmail)
	}

	checkUser, err = r.repo.GetUserByNickname(ctx, user.Nickname)
	if err != nil {
		errsList = append(errsList, errs.ErrDBInternal)
		return errsList
	}
	if checkUser != emptyUser {
		errsList = append(errsList, errs.ErrDBUniqueNickname)
	}

	if len(errsList) > 0 {
		return errsList
	}

	err = r.repo.RegisterUser(ctx, user)
	if err != nil {
		errsList = append(errsList, errs.ErrDBInternal)
	}

	return errsList
}

func (r *RepositoryService) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if user.Nickname != "" {
		checkUser, err := r.repo.GetUserByNickname(ctx, user.Nickname)
		if err != nil {
			return emptyUser, errs.ErrDBInternal
		}
		if checkUser != emptyUser && checkUser.UserID != user.UserID {
			return emptyUser, errs.ErrDBUniqueNickname
		}
	}

	err := r.repo.UpdateUser(ctx, user)
	if err != nil {
		return emptyUser, errs.ErrDBInternal
	}

	updatedUser, err := r.repo.GetUserById(ctx, user.UserID)
	if err != nil {
		return emptyUser, errs.ErrDBInternal
	}

	return updatedUser, nil
}
