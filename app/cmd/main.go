package main

import (
	"net/http"
	_ "net/http/pprof"
	"yt-solutions-soft/internal/service"
	"yt-solutions-soft/internal/websocket"
)

func main() {
	go func() {
		http.ListenAndServe(":8088", nil)
	}()

	s := service.NewService()
	a := websocket.NewAdapter(s.YouTubeChecker, s.ProxyChecker)

	server := websocket.NewServer(a.InitRoutes())

	if err := server.Start(); err != nil {
		panic(err)
	}

}
