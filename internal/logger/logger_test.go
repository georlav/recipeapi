package logger_test

import (
	"log"
	"os"
	"testing"

	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/logger"
)

func TestNew(t *testing.T) {
	cfg := config.Logger{
		LogLevel:     1,
		EnableStdout: false,
	}

	t.Run("Should initialize a new logger", func(t *testing.T) {
		l := logger.NewLogger(cfg)
		if l.GetLevel() != 1 {
			log.Fatalf("Invalid level, expected %d got %d", cfg.LogLevel, 1)
		}
	})

	t.Run("Should initialize and inject a custom hook", func(t *testing.T) {
		h := graylog.NewGraylogHook(
			"127.0.0.1:1234",
			map[string]interface{}{
				"application": "testid",
				"environment": "testenv",
				"version":     "1",
			},
		)

		noOpts := logger.NewLogger(cfg)
		withOpts := logger.NewLogger(cfg, logger.SetHook(h))

		if len(noOpts.Hooks) == len(withOpts.Hooks) {
			log.Fatal("Invalid number of hooks")
		}
	})

	t.Run("Should initialize and change standard output", func(t *testing.T) {
		l := logger.NewLogger(cfg, logger.SetOutput(os.Stdout))
		if l.Out != os.Stdout {
			log.Fatal("Invalid output")
		}
	})
}
