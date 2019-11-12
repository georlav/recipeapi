package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config object
type Config struct {
	APP    APP    `json:"app"`
	Server Server `json:"Server"`
}

// APP holds general app configuration values
type APP struct {
	Debug bool `json:"debug"`
}

// Server object holds the base configuration for the http server
// ReadTimeout is the maximum duration for reading the entire request, including the body (seconds)
// WriteTimeout is the maximum duration before timing out writes the response (seconds)
// IdleTimeout is the maximum amount of time to wait for the next request when keep-alive is enabled (seconds)
type Server struct {
	Protocol     string
	Port         int
	ReadTimeout  int64
	WriteTimeout int64
	IdleTimeout  int64
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
