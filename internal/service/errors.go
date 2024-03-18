package service

import "errors"

// штуку с ошибками надо решить: если обращаться к ним через handler,
// то ошибка import cycle not allowed при запуска,
// поэтому пришлось тупо скопировать пока
var (
	ErrDBUniqueEmail    = errors.New("user with this email already exists (can't register)")
	ErrDBUniqueNickname = errors.New("user with this nickname already exists (can't register)")
	ErrDBInternal       = errors.New("internal db error")
)
