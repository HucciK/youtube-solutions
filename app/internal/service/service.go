package service

import (
	"yt-solutions-soft/internal/models"
	"yt-solutions-soft/internal/service/checker"
	"yt-solutions-soft/internal/service/proxy"
)

type YouTubeChecker interface {
	CheckFolder(checkPath, savePath string, updatesChan chan models.Update, proxies []models.Proxy)
}

type ProxyChecker interface {
	SpeedTest(proxies []models.Proxy, updatesChan chan models.Update)
}

type Service struct {
	YouTubeChecker
	ProxyChecker
}

func NewService() *Service {
	return &Service{
		YouTubeChecker: checker.NewYouTubeChekcer(),
		ProxyChecker:   proxy.NewProxyChecker(),
	}
}
