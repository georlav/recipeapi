package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config object
type Config struct {
	APP      APP
	Server   Server
	Database Database
	Logger   Logger
	Token    Token
}

// APP holds general app configuration values
type APP struct {
	Version string
}

// Server holds the base configuration for the http server
// ReadTimeout is the maximum duration for reading the entire request, including the body (seconds)
// WriteTimeout is the maximum duration before timing out writes the response (seconds)
// IdleTimeout is the maximum amount of time to wait for the next request when keep-alive is enabled (seconds)
type Server struct {
	Scheme       string
	Host         string
	Port         int
	ReadTimeout  int64
	WriteTimeout int64
	IdleTimeout  int64
}

// Database holds the base configuration for the application db storage
type Database struct {
	Host         string
	Port         int16
	Username     string
	Password     string
	Database     string
	MaxIdleConns int
	MaxOpenConns int
}

// Logger holds the configuration for logging
type Logger struct {
	LogLevel     uint8
	EnableStdout bool
	ReportCaller bool
	AppVersion   string
}

// Token holds configuration for tokens
type Token struct {
	Secret string
	TTL    int64 // Minutes
}

// New returns a new config, by default it looks for config files in the current working directory, if your config
// is locate somewhere path the path as second argument
func New(name string, path ...string) (*Config, error) {
	v := viper.New()
	// Setup file to read
	paths := append([]string{"."}, path...)
	for i := range paths {
		v.AddConfigPath(paths[i])
	}
	v.SetConfigName(name)

	// Read ENV variables with recipe prefix
	v.SetEnvPrefix("recipe")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config, %w", err)
	}

	// Override values using cli arguments
	setFlags(v)

	c := Config{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config, %w", err)
	}
	return &c, nil
}

// Create flags for available options (integers, strings, booleans), use cli arguments to override config values
func setFlags(v *viper.Viper) {
	keys := v.AllKeys()
	for i := range keys {
		k := v.Get(keys[i])
		switch k.(type) {
		case int, int8, int16, int32:
			v.Set(
				keys[i],
				*flag.Int(keys[i], k.(int), fmt.Sprintf("Sets value for %s", keys[i])),
			)
		case int64:
			v.Set(
				keys[i],
				*flag.Int64(keys[i], k.(int64), fmt.Sprintf("Sets value for %s", keys[i])),
			)
		case string:
			v.Set(
				keys[i],
				flag.String(keys[i], k.(string), fmt.Sprintf("Sets value for %s", keys[i])),
			)
		case bool:
			v.Set(
				keys[i],
				*flag.Bool(keys[i], k.(bool), fmt.Sprintf("Sets value for %s", keys[i])),
			)
		}
	}

	flag.Parse()
}
