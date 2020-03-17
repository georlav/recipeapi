package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config object
type Config struct {
	APP      APP      `json:"app"`
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
	Token    Token    `json:"token"`
}

// APP holds general app configuration values
type APP struct {
	Version string `json:"version"`
}

// Server holds the base configuration for the http server
// ReadTimeout is the maximum duration for reading the entire request, including the body (seconds)
// WriteTimeout is the maximum duration before timing out writes the response (seconds)
// IdleTimeout is the maximum amount of time to wait for the next request when keep-alive is enabled (seconds)
type Server struct {
	Scheme       string `json:"scheme"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int64  `json:"readTimeout"`
	WriteTimeout int64  `json:"writeTimeout"`
	IdleTimeout  int64  `json:"idleTimeout"`
}

// Database holds the base configuration for the application db storage
type Database struct {
	Host         string `json:"host"`
	Port         int16  `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

// Logger holds the configuration for logging
type Logger struct {
	LogLevel     uint8  `json:"logLevel"`
	EnableStdout bool   `json:"enableStdout"`
	ReportCaller bool   `json:"reportCaller"`
	AppVersion   string `json:"appVersion"`
}

// Token holds configuration for tokens
type Token struct {
	Secret string `json:"secret"`
	TTL    int64  `json:"ttl"` // Minutes
}

// Load loads a json config file and returns a config object
func Load(cfgFile string) (*Config, error) {
	b, err := ioutil.ReadFile("./" + cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s, %w", cfgFile, err)
	}

	cfg := Config{}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s, %w", cfgFile, err)
	}
	cfg.Logger.AppVersion = cfg.APP.Version

	return &cfg, nil
}
