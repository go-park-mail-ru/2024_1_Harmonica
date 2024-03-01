package handler

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
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]*$", nickname)
	if !match {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	return len(password) > 8 && len(password) < 20
}
