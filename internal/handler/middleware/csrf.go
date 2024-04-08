package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRF(next http.Handler) http.Handler {
	return csrf.Protect([]byte("32-byte-long-auth-key"),
		csrf.SameSite(csrf.SameSiteNoneMode),
		csrf.HttpOnly(true))(next)
}
