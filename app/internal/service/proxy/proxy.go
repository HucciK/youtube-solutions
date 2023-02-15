package proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"yt-solutions-soft/internal/models"
)

type ProxyChecker struct {
	Client *http.Client
}

func NewProxyChecker() *ProxyChecker {
	return &ProxyChecker{
		Client: &http.Client{},
	}
}

func (p *ProxyChecker) SpeedTest(proxies []models.Proxy, updatesChan chan models.Update) {

	if len(proxies) == 0 {
		return
	}

	for _, proxy := range proxies {
		var upd models.Update

		if err := p.setNewProxy(proxy); err != nil {
			upd.Err = err
			updatesChan <- upd
			continue
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*2000)

		req, err := http.NewRequestWithContext(ctx, "GET", "https://youtube.com", nil)
		if err != nil {
			upd.Err = err
			updatesChan <- upd
			continue
		}

		start := time.Now()
		res, err := p.Client.Do(req)
		if err != nil {
			upd.Err = err
			updatesChan <- upd
			continue
		}

		if res.StatusCode != 200 {
			upd.Err = errors.New("proxy doesn't response")
			updatesChan <- upd
			continue
		}

		upd.Data = int(time.Now().Sub(start).Milliseconds())
		updatesChan <- upd
	}
	close(updatesChan)
}

func (p *ProxyChecker) setNewProxy(proxy models.Proxy) error {
	proxyStr := fmt.Sprintf("http://%s:%s@%s:%s", proxy.User, proxy.Password, proxy.IP, proxy.Port)
	proxyUrl, err := url.Parse(proxyStr)
	if err != nil {
		return err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}

	p.Client.Transport = transport

	return nil
}
