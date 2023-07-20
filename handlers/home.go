package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kisstc/go-backend-template/server"
)

// se encargara de procesar la peticion de la ruta principal /
// los mensajes que le enviaremos al cliente
type HomeResponse struct {
	Message string `json:"message"` // serializamos a json, tomara el nombre en minuscula
	Status  bool   `json:"status"`  // serializamos a json
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to Ganadores",
			Status:  true,
		})
	}
}
