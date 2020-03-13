package logger

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/sirupsen/logrus"
)

type Option func(*Logger)

// Hook inject your hook integration into logrus
func SetHook(hook logrus.Hook) Option {
	return func(args *Logger) {
		args.AddHook(hook)
	}
}

// SetOutput to stdOut or to any custom buffer
func SetOutput(out io.Writer) Option {
	return func(args *Logger) {
		args.SetOutput(out)
	}
}

type Logger struct {
	*logrus.Logger
	cfg config.Logger
}

// NewLogger creates, initializes and returns a new logger
func NewLogger(cfg config.Logger, options ...Option) *Logger {
	l := Logger{
		Logger: logrus.New(),
		cfg:    cfg,
	}

	l.SetReportCaller(cfg.ReportCaller)
	l.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	l.SetLevel(cfg.LogLevel)
	if !cfg.APP.Debug {
		l.SetOutput(ioutil.Discard)
	}

	for i := range options {
		options[i](&l)
	}

	return &l
}
