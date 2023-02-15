package main

import (
	"fmt"
	"yt-solutions-telegram-dashboard/config"
	"yt-solutions-telegram-dashboard/internal/service"
	"yt-solutions-telegram-dashboard/internal/transport"
	"yt-solutions-telegram-dashboard/internal/transport/handler"
)

const (
	clean = 86400
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	cli := service.NewClient(cfg.Host, cfg.Token, cfg.Backend)
	sm := service.NewStateMachine()
	tm := service.NewTransactionManager(clean)
	c := service.NewCallbackService(cli, sm, cfg.DownloadLink, cfg.GuideLink)
	m := service.NewMessageService(cli, sm, tm)

	h := handler.NewHandler(c, m)

	server, cert := transport.NewServer(cfg.ServerConfig, h.InitRoutes())

	if err = server.Start(cert); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
