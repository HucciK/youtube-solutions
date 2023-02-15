package websocket

import "yt-solutions-soft/internal/models"

type CheckerDataReceive struct {
	Paths   CheckerPaths   `json:"paths"`
	Proxies []models.Proxy `json:"proxy"`
}

type CheckerPaths struct {
	CheckPath string `json:"check_path"`
	SavePath  string `json:"save_path"`
}

type Proxies struct {
	ProxyList []models.Proxy `json:"proxy_list"`
}
