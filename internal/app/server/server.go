package server

import (
	"log"
	"net/http"

	"github.com/Akezhan1/forum/internal/app/repository"
	"github.com/Akezhan1/forum/internal/app/service"

	"github.com/Akezhan1/forum/internal/app/handler"
)

type Server struct {
	httpServer *http.Server
}

func New(config *Config) *Server {
	db, err := repository.OpenDB(config.DBDriver, config.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)

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
