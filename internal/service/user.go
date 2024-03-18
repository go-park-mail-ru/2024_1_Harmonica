package service

import (
	"harmonica/internal/entity"
	"harmonica/internal/repository"
)

//// штуку с ошибками надо решить: если обращаться к ним через handler,
//// то ошибка import cycle not allowed при запуска,
//// поэтому пришлось тупо скопировать пока
//var (
//	ErrDBUniqueEmail    = errors.New("user with this email already exists (can't register)")
//	ErrDBUniqueNickname = errors.New("user with this nickname already exists (can't register)")
//	ErrDBInternal       = errors.New("internal db error")
//)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(r *repository.Repository) *UserService {
	return &UserService{repo: r}
}

// регистрация: SELECT с поиском пользователя здесь, а не из хэндлера

func (u UserService) GetUserByEmail(email string) (entity.User, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return entity.User{}, ErrDBInternal
	}
	return user, nil
}

func (u UserService) GetUserByNickname(nickname string) (entity.User, error) {
	user, err := u.repo.GetUserByNickname(nickname)
	if err != nil {
		return entity.User{}, ErrDBInternal
	}
	return user, nil
}

func (u UserService) GetUserById(id int64) (entity.User, error) {
	user, err := u.repo.GetUserById(id)
	if err != nil {
		return entity.User{}, ErrDBInternal
	}
	return user, nil
}

func (u UserService) RegisterUser(user entity.User) []error {
	isEmailUnique, isNicknameUnique := false, false
	var errs []error

	// Checking for unique fields
	user, err := u.repo.GetUserByEmail(user.Email)
	if err != nil {
		errs = append(errs, ErrDBInternal)
		return errs
	}
	if user == (entity.User{}) {
		isEmailUnique = true
	}
	user, err = u.repo.GetUserByNickname(user.Nickname)
	if err != nil {
		errs = append(errs, ErrDBInternal)
		return errs
	}
	if user == (entity.User{}) {
		isNicknameUnique = true
	}

	if !isEmailUnique {
		errs = append(errs, ErrDBUniqueEmail)
	}
	if !isNicknameUnique {
		errs = append(errs, ErrDBUniqueNickname)
	}
	if len(errs) > 0 {
		return errs
	}

	err = u.repo.RegisterUser(user)
	if err != nil {
		errs = append(errs, ErrDBInternal)
	}

	return errs
}
