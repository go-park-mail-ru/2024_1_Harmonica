package utils

import (
	"net/mail"
	"regexp"
	"unicode"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateNickname(nickname string) bool {
	// можно - цифры, латинские буквы и (!!!) знак нижнего подчеркивания _
	// обязательно - длина от 3 до 20
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", nickname)
	return match
}

func ValidatePassword(password string) bool {
	// можно - цифры и латинские буквы
	// обязательно - длина от 8 до 24, наличие хотя бы 1 заглавной буквы, наличие хотя бы 1 цифры
	if len(password) < 8 || len(password) > 24 {
		return false
	}
	hasUppercase, hasDigit := false, false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
		if unicode.IsLetter(char) {
			if !unicode.Is(unicode.Latin, char) {
				return false
			}
		}
	}
	return hasUppercase && hasDigit
}
