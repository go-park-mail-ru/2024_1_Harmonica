package utils

import (
	"net/mail"
	"regexp"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateNickname(nickname string) bool {
	if len(nickname) < 3 || len(nickname) > 20 {
		return false
	}
	// вроде так:
	// можно - цифры, латинские буквы и (!!!) знак нижнего подчеркивания _
	// обязательно - длина от 3 до 20
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]*$", nickname)
	if !match {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 15 {
		return false
	}
	// вроде так:
	// можно - цифры и латинские буквы
	// обязательно - длина от 8 до 15, наличие хотя бы 1 заглавной буквы
	//match, _ := regexp.MatchString("^(?=.*[A-Z])[A-Za-z0-9]*$", password)
	match, _ := regexp.MatchString("^[A-Za-z0-9]*$", password) // это мое для тестов
	if !match {
		return false
	}
	return true
}
