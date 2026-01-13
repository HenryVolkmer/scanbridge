package main

import (
	"encoding/json"
	"net/url"
	"os"
)

type Config struct {
	IsAutodiscovery bool `json:"isAutodiscovery"`
	Devices []*ScanDevice `json:"devices"`
	Smtp *SmtpConfig `json:"smtp"`
	IsDebug bool
}

type SmtpConfig struct {
	Host *url.URL `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Sender string `json:"sender"`
	Recipient string `json:"recipient"`
	Subject string `json:"subject"`
}

// New unmarshals the given config Filename into
// the Config struct
func NewConfig(cfgFileName string) (*Config, error) {
	cfgBytes, err := os.ReadFile(cfgFileName)
	if err != nil {
		return  nil, err
	}
	cfg := &Config{}
	json.Unmarshal(cfgBytes, cfg)
	return cfg, nil
}