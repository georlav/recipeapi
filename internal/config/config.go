package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// Config object
type Config struct {
	APP      APP      `json:"app"`
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

// APP holds general app configuration values
type APP struct {
	Version int  `json:"version"`
	Debug   bool `json:"debug"`
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
	Host     string
	Port     int16
	Username string
	Password string
	Database string
}

// Logger holds the configuration for logging
type Logger struct {
	LogLevel     logrus.Level
	EnableStdout bool
	ReportCaller bool
	APP          APP
}

// Load loads a json config file and returns a config object
func Load(cfgFile string) (cfg *Config, err error) {
	b, err := ioutil.ReadFile("./" + cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s, %w", cfgFile, err)
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s, %w", cfgFile, err)
	}

	return cfg, nil
}
