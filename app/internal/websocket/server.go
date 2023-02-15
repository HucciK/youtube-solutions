package websocket

import (
	"net/http"
	"time"
)

type Server struct {
	Server *http.Server
}

func NewServer(h http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         ":8080",
			Handler:      h,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
	}
}

func (s Server) Start() error {
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
