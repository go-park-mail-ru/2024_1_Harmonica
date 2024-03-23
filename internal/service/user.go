package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errors_list"
)

var emptyUser = entity.User{}

func (r RepositoryService) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := r.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return emptyUser, errors_list.ErrDBInternal
	}
	return user, nil
}

func (r RepositoryService) GetUserByNickname(ctx context.Context, nickname string) (entity.User, error) {
	user, err := r.repo.GetUserByNickname(ctx, nickname)
	if err != nil {
		return emptyUser, errors_list.ErrDBInternal
	}
	return user, nil
}

func (r RepositoryService) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	user, err := r.repo.GetUserById(ctx, id)
	if err != nil {
		return emptyUser, errors_list.ErrDBInternal
	}
	return user, nil
}

func (r RepositoryService) RegisterUser(ctx context.Context, user entity.User) []error {
	isEmailUnique, isNicknameUnique := false, false
	var errs []error

	// Checking for unique fields
	checkUser, err := r.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		errs = append(errs, errors_list.ErrDBInternal)
		return errs
	}
	if checkUser == emptyUser {
		isEmailUnique = true
	}
	checkUser, err = r.repo.GetUserByNickname(ctx, user.Nickname)
	if err != nil {
		errs = append(errs, errors_list.ErrDBInternal)
		return errs
	}
	if checkUser == emptyUser {
		isNicknameUnique = true
	}

	if !isEmailUnique {
		errs = append(errs, errors_list.ErrDBUniqueEmail)
	}
	if !isNicknameUnique {
		errs = append(errs, errors_list.ErrDBUniqueNickname)
	}
	if len(errs) > 0 {
		return errs
	}

	err = r.repo.RegisterUser(ctx, user)
	if err != nil {
		errs = append(errs, errors_list.ErrDBInternal)
	}

	return errs
}

func (r RepositoryService) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	// Checking for unique field
	checkUser, err := r.repo.GetUserByNickname(ctx, user.Nickname)
	if err != nil {
		return emptyUser, errors_list.ErrDBInternal
	}
	if checkUser != emptyUser {
		return emptyUser, errors_list.ErrDBUniqueNickname
	}

	err = r.repo.UpdateUser(ctx, user)
	if err != nil {
		return emptyUser, errors_list.ErrDBInternal
	}

	updatedUser, err := r.repo.GetUserByNickname(ctx, user.Nickname)
	if err != nil {
		return emptyUser, errors_list.ErrDBInternal
	}

	return updatedUser, nil
}
