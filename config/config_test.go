package config_test

import (
	"strings"
	"testing"

	"github.com/georlav/recipeapi/config"
)

func TestLoad(t *testing.T) {
	t.Run("Should successfully parse", func(t *testing.T) {
		_, err := config.Load("testdata/valid.json")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Should fail to parse due to invalid format", func(t *testing.T) {
		_, err := config.Load("testdata/invalid.json")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "failed to unmarshal") {
			t.Fatal("unexpected error")
		}
	})

	t.Run("Should fail to open config file", func(t *testing.T) {
		_, err := config.Load("testdata/notfound.json")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "failed to open") {
			t.Fatal("unexpected error")
		}
	})
}
