package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerConfig `json:"server"`
	DBСonfig     `json:"db"`
	OrdersConfig `json:"orders"`
	Version      string `json:"version"`
}

type ServerConfig struct {
	Addr string `json:"addr"`
}

type DBСonfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	Sslmode  string `json:"sslmode"`
}

type OrdersConfig struct {
	MaxFree     int `json:"max_free"`
	MaxLifetime int `json:"max_lifetime"`
	Price       int `json:"price"`
	Renewal     int `json:"renewal"`
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

	return &config, nil
}

func (db DBСonfig) ToDataSource() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", db.Host, db.Port, db.User, db.DBName, db.Password, db.Sslmode)
}
