package middleware

import (
	"context"
	"log"
	"net/http"
	"rest-ws/models"
	"rest-ws/server"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var (
	NO_AUTH_NEEDED = []string{
		"/", "/health", "/signup", "/login",
	}
)

func shouldCheckToken(path string) bool {
	for _, v := range NO_AUTH_NEEDED {
		if v == path {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(s server.Server) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(s.Config().JWTSecret), nil
				})
			if err != nil || !token.Valid {
				log.Printf("error parsing token: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userToken", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
