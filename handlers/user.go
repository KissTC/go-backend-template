package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kisstc/go-backend-template/models"
	"github.com/kisstc/go-backend-template/repository"
	"github.com/kisstc/go-backend-template/server"
	"github.com/segmentio/ksuid"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
		// se rellena de informacion el struct aquí
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			// se envia codigo 400
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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
			Password: request.Password,
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
