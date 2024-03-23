package middleware

//
//import (
//	"errors"
//	"fmt"
//	"net/http"
//	"strings"
//)
//
//func AuthMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		c, err := r.Cookie("session_token")
//		if err != nil {
//			fmt.Println("no auth at", r.URL.Path)
//			http.Redirect(w, r, "/", http.StatusUnauthorized)
//			return
//		}
//
//		authHeader := c.Request.Header.Get("Authorization")
//		t := strings.Split(authHeader, " ")
//		if len(t) == 2 {
//			authToken := t[1]
//			authorized, err := tokenutil.IsAuthorized(authToken, secret)
//			if authorized {
//				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
//				if err != nil {
//					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
//					c.Abort()
//					return
//				}
//				c.Set("x-user-id", userID)
//				c.Next()
//				return
//			}
//			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
//			c.Abort()
//			return
//		}
//		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
//		c.Abort()
//	}
//}
//
//func CheckAuth(r *http.Request) (string, int64, error) {
//	c, err := r.Cookie("session_token")
//	if err != nil {
//		if errors.Is(err, http.ErrNoCookie) {
//			return "", 0, nil
//		}
//		return "", 0, nil
//	}
//	sessionToken := c.Value
//	s, exists := sessions.Load(sessionToken)
//	if !exists {
//		return "", 0, nil
//	}
//	if s.(Session).IsExpired() {
//		sessions.Delete(sessionToken)
//		return "", 0, nil
//	}
//	return sessionToken, s.(Session).UserId, nil
//	// в мидлваре прокинуть session_token в контекст, чтобы он был досупен далее в ручке
//}
