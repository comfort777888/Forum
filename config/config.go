package config

import (
	"encoding/json"
	"os"
	"time"
)

const (
	filename   = "config/config.json"
	maxHeader  = 1 >> 20
	writeTO    = 5 * time.Second
	shutDownTO = 3 * time.Second
)

type Config struct {
	Port            string `json:"port"`
	Host            string `json:"host"`
	MaxHeaderBytes  int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeOut time.Duration
	DbNameAndPath   string `json:"dbNameAndPath"`
	DbDriver        string `json:"dbDriver"`
	CtxTimeout      int    `json:"ctxTimeout"`
}

func New() (*Config, error) {
	var config *Config
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return nil, err
	}
	// defer configFile.Close()

	return &Config{
		Port:            config.Port,
		MaxHeaderBytes:  maxHeader,
		ReadTimeout:     writeTO,
		WriteTimeout:    writeTO,
		ShutdownTimeOut: shutDownTO,
		DbNameAndPath:   config.DbNameAndPath,
		DbDriver:        config.DbDriver,
		CtxTimeout:      config.CtxTimeout,
	}, nil
}
