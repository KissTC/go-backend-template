package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/kisstc/go-backend-template/models"
	"github.com/kisstc/go-backend-template/server"
)

// rutas no necesarias para revisar
var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signip",
	}
)

func shouldCheckToken(path string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(path, p) {
			return false
		}
	}

	return true
}

// revisa si la ruta debe checarse o no. Si no, avanza
// si debe checarse, usamos el token para ver la autorizacion del header authorization
// si el token está bien el usuario se valida y se avanza, si no, enviamos unauthorized
func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			// el token viene del header de authorization
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenString, models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
