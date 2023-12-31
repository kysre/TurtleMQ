package main

import (
	"github.com/sirupsen/logrus"
)

// LoggingConfig the logger's configuration structure
type LoggingConfig struct {
	Level string
}

// ConfigureLogging handler config logger based on the given configuration
func provideLogger(config *Config) (*logrus.Logger, error) {
	logger := logrus.New()
	if config.Logging.Level != "" {
		level, err := logrus.ParseLevel(config.Logging.Level)
		if err != nil {
			return nil, err
		}
		logger.SetLevel(level)
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})

	return logger, nil
}
