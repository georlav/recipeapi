package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// Config object
type Config struct {
	APP    APP    `json:"app"`
	Server Server `json:"server"`
	Mongo  Mongo  `json:"mongo"`
}

// APP holds general app configuration values
type APP struct {
	Debug bool `json:"debug"`
}

// Mongo holds the configuration for mongo database
type Mongo struct {
	Host                      string `json:"host"`
	Port                      int16  `json:"port"`
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	PoolSize                  uint16 `json:"poolSize"`
	Timeout                   time.Duration
	SetServerSelectionTimeout time.Duration
	SetConnectTimeout         time.Duration
	SetSocketTimeout          time.Duration
	SetMaxConnIdleTime        time.Duration
	SetRetryWrites            bool   `json:"setRetryWrites"`
	Database                  string `json:"database"`
	RecipeCollection          string `json:"recipeCollection"`
}

// Server object holds the base configuration for the http server
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
