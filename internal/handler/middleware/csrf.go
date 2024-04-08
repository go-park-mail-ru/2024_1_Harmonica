package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRF(next http.Handler) http.Handler {
	return csrf.Protect([]byte("32-byte-long-auth-key"),
		csrf.HttpOnly(true))(next)
}

//csrf.SameSite(csrf.SameSiteNoneMode)
