package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kisstc/go-backend-template/models"
	"github.com/kisstc/go-backend-template/repository"
	"github.com/kisstc/go-backend-template/server"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

//tipo nuevo para procesar el request de autenticacion
type SignUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{}
		// se rellena de informacion el struct aquí
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			// se envia codigo 400
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// usando bcrypt para hacer el hash del password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// libreria para generar id aleatorio
		id, err := ksuid.NewRandom()
		if err != nil {
			// se envia codigo 500, error del servidor
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// si todo esta listo creamos la variable de usuario para asignar
		var user = models.User{
			Email:    request.Email,
			Password: string(hashedPassword), // enviamos el password hashed, convertido a string
			Id:       id.String(),
		}

		// llamamos a repository
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			// se envia codigo 500, error del servidor
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// si se envian bien los valores usamos el repositorio
		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if user == nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized) // 401
		}

		// compara el password de la bd con el del request
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized) // 401
			return
		}

		// decodigicar un nuevo token de respuesta, hasta este punto sabemos que el usuario existe
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(), // 2 días de token valido
			},
		}

		// ya teniendo los claims generamos el nuevo token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// token creado, ahora necesitamos firmarlo
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//no hay errores a este punto, respondemos al cliente con el nuevo response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}
