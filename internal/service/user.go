package service

import (
	"harmonica/internal/entity"
	"harmonica/internal/entity/errors_list"
)

//в будущем разделить ошибки на те, что внутри сервиса и те, что идут к пользователю
//ошибки базы
//ошибки сервиса
//ошибки для пользователя

//type UserService struct {
//	repo *repository.Repository
//}
//
//func NewUserService(r *repository.Repository) *UserService {
//	return &UserService{repo: r}
//}

func (r RepositoryService) GetUserByEmail(email string) (entity.User, error) {
	user, err := r.repo.GetUserByEmail(email)
	if err != nil {
		return entity.User{}, errors_list.ErrDBInternal
	}
	return user, nil
}

func (r RepositoryService) GetUserByNickname(nickname string) (entity.User, error) {
	user, err := r.repo.GetUserByNickname(nickname)
	if err != nil {
		return entity.User{}, errors_list.ErrDBInternal
	}
	return user, nil
}

func (r RepositoryService) GetUserById(id int64) (entity.User, error) {
	user, err := r.repo.GetUserById(id)
	if err != nil {
		return entity.User{}, errors_list.ErrDBInternal
	}
	return user, nil
}

func (r RepositoryService) RegisterUser(user entity.User) []error {
	isEmailUnique, isNicknameUnique := false, false
	var errs []error

	// Checking for unique fields
	checkUser, err := r.repo.GetUserByEmail(user.Email)
	if err != nil {
		errs = append(errs, errors_list.ErrDBInternal)
		return errs
	}
	if checkUser == (entity.User{}) {
		isEmailUnique = true
	}
	checkUser, err = r.repo.GetUserByNickname(user.Nickname)
	if err != nil {
		errs = append(errs, errors_list.ErrDBInternal)
		return errs
	}
	if checkUser == (entity.User{}) {
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

	err = r.repo.RegisterUser(user)
	if err != nil {
		errs = append(errs, errors_list.ErrDBInternal)
	}

	return errs
}
