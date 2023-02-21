package transport

import (
	"net/http"
	"yt-solutions-server/config"
)

type Server struct {
	Server http.Server
}

func NewServer(cfg config.ServerConfig, h http.Handler) *Server {

	//APP_IP := os.Getenv("APP_IP")
	//APP_PORT := os.Getenv("APP_PORT")

	return &Server{
		Server: http.Server{
			//Addr:    APP_IP + ":" + APP_PORT,
			Addr:    ":8080",
			Handler: h,
		},
	}
}

func (s Server) Start() error {
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
