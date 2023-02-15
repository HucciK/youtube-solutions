package main

import (
	"yt-solutions-server/config"
	"yt-solutions-server/internal/repository"
	"yt-solutions-server/internal/services"
	"yt-solutions-server/internal/transport"
	"yt-solutions-server/internal/transport/handler"
)

const token = "5678637087:AAH34YJMUi3NMf55gGTlV6cmZgJF4TzkXMo"

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := repository.NewPostgres(cfg.DBСonfig)
	if err != nil {
		panic(err)
	}

	urepo := repository.NewUserRepo(db)
	krepo := repository.NewKeyRepo(db)
	orepo := repository.NewOrderRepo(db)

	uservice := services.NewUserService(urepo, token)
	kservice := services.NewKeyService(krepo)
	oservice := services.NewOrderService(urepo, krepo, orepo, cfg.OrdersConfig)

	h := handler.NewHandler(uservice, kservice, oservice, cfg.Version)

	s := transport.NewServer(cfg.ServerConfig, h.InitRoutes())

	if err = s.Start(); err != nil {
		panic(err)
	}
}
