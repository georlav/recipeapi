package config_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/spf13/viper"
)

func TestLoad(t *testing.T) {
	t.Run("Should successfully parse", func(t *testing.T) {
		_, err := config.New("valid", "testdata")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Should fail to parse due to invalid format", func(t *testing.T) {
		_, err := config.New("invalid", "testdata")
		if err == nil {
			t.Fatal("Expected to have ConfigParseError")
		}

		if !errors.As(err, &viper.ConfigParseError{}) {
			t.Fatalf("unexpected error, %s, type: %s", errors.Unwrap(err), reflect.TypeOf(errors.Unwrap(err)))
		}
	})

	t.Run("Should fail to open config file", func(t *testing.T) {
		_, err := config.New("notfound", "testdata")
		if err == nil {
			t.Fatal("Expected to have ConfigFileNotFoundError")
		}
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			t.Fatalf("unexpected error, %s, type: %s", errors.Unwrap(err), reflect.TypeOf(errors.Unwrap(err)))
		}
	})

	t.Run("Should override config value with environmental value", func(t *testing.T) {
		os.Setenv("RECIPE_TOKEN_TTL", "66")
		defer os.Unsetenv("RECIPE_TOKEN_TTL")

		cfg, err := config.New("valid", "testdata")
		if err != nil {
			t.Fatal(err)
		}

		if cfg.Token.TTL != 66 {
			t.Fatalf("Token TTL expected to have value %d got %d", 66, cfg.Token.TTL)
		}
	})
}
