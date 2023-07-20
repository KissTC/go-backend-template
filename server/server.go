package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// configuracion que nuestro servidor requiere para conectarse
type Config struct {
	Port        string
	JWTSecret   string // clave para generar tokens
	DatabaseUrl string // conexion a bd
}

type Server interface {
	Config() *Config
}

// manejador de servidores
type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("Secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("Database url is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	// b.router = mux.NewRouter() no necesario?
	binder(b, b.router)
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
