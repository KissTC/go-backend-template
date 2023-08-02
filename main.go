package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kisstc/go-backend-template/handlers"
	"github.com/kisstc/go-backend-template/server"
)

func main() {
	// cargar variables de entorno
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	// creamos servidor nuevo
	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})
	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)

}

// funcion bind para start server
func BindRoutes(s server.Server, r *mux.Router) {
	// .Methods() asignamos el metodo http que usará
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)          //usará GET
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost) // usara POST

}
