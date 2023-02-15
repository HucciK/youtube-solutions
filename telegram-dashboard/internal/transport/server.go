package transport

import (
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"yt-solutions-telegram-dashboard/config"
)

type Server struct {
	Server http.Server
}

func NewServer(cfg config.ServerConfig, h http.Handler) (*Server, *autocert.Manager) {

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
	}

	APP_IP := os.Getenv("APP_IP")
	APP_PORT := os.Getenv("APP_PORT")

	fmt.Println(APP_IP, APP_PORT)
	return &Server{
		Server: http.Server{
			Addr:    APP_IP + ":" + APP_PORT,
			Handler: h,
		},
	}, &certManager
}

func (s Server) Start(certManager *autocert.Manager) error {

	//go http.ListenAndServe(APP_IP+":"+APP_PORT, certManager.HTTPHandler(nil))

	//if err := s.Server.ListenAndServeTLS("", ""); err != nil {
	//	return err
	//}

	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
