package websocket

import (
	"fmt"
	"net/http"
	"yt-solutions-soft/internal/models"

	"golang.org/x/net/websocket"
)

type YouTubeChecker interface {
	CheckFolder(checkPath, savePath string, updatesChan chan models.Update, proxies []models.Proxy)
}

type ProxyChecker interface {
	SpeedTest(proxies []models.Proxy, updatesChan chan models.Update)
}

type Adapter struct {
	YouTubeChecker
	ProxyChecker
}

func NewAdapter(ytC YouTubeChecker, p ProxyChecker) *Adapter {
	return &Adapter{
		YouTubeChecker: ytC,
		ProxyChecker:   p,
	}
}

func (a Adapter) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/checker", websocket.Handler(a.CheckerCall))
	mux.Handle("/proxy_checker", websocket.Handler(a.ProxyCheckerCall))

	return mux
}

func (a Adapter) CheckerCall(ws *websocket.Conn) {
	fmt.Println(ws.Request().Header.Get("Sec-Websocket-Protocol"))

	var data CheckerDataReceive

	if err := websocket.JSON.Receive(ws, &data); err != nil {
		return
	}

	updatesChan := make(chan models.Update)

	go a.YouTubeChecker.CheckFolder(data.Paths.CheckPath, data.Paths.SavePath, updatesChan, data.Proxies)

	for update := range updatesChan {
		websocket.JSON.Send(ws, update)
	}

}

func (a Adapter) ProxyCheckerCall(ws *websocket.Conn) {
	var data Proxies

	if err := websocket.JSON.Receive(ws, &data); err != nil {
		return
	}

	updatesChan := make(chan models.Update, len(data.ProxyList))
	go a.ProxyChecker.SpeedTest(data.ProxyList, updatesChan)

	for update := range updatesChan {
		websocket.JSON.Send(ws, update)
	}

}
