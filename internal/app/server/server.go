package server

import (
	"log"
	"net/http"

	"github.com/Akezhan1/forum/internal/app/handler"
)

type Server struct {
	httpServer *http.Server
}

func New(config *Config) *Server {
	handler := handler.NewHandler()
	return &Server{
		httpServer: &http.Server{
			Addr:    config.Addr,
			Handler: handler.InitRouter(),
		},
	}
}

func (s *Server) Start() error {
	log.Println("starting api server")

	return s.httpServer.ListenAndServe()
}
