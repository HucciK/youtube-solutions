package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerConfig   `json:"server"`
	TelegramConfig `json:"telegram"`
	DownloadLink   string `json:"downloadLink"`
	GuideLink      string `json:"guideLink"`
}

type ServerConfig struct {
	Addr string `json:"addr"`
}

type TelegramConfig struct {
	Token   string `json:"token"`
	Host    string `json:"host"`
	Backend string `json:"backend"`
}

func NewConfig() (*Config, error) {
	var config Config

	data, err := os.ReadFile("../config/config.json")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, err
}
