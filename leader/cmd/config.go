package main

import (
	"strings"

	"github.com/kysre/TurtleMQ/leader/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config the application's configuration structure
type Config struct {
	ConfigFile string
	Logging    LoggingConfig
	Queue      QueueConfig
	Leader     LeaderConfig
	Metric     MetricConfig
}

type QueueConfig struct {
	ListenPort int
}

type LeaderConfig struct {
	ListenPort  int
	ReplicaHost string

	DataNodeStateCheckPeriod     int
	DataNodeRemainingCheckPeriod int
}

type MetricConfig struct {
	ListenPort int
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	// Read Config from ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("LEADER")
	viper.AutomaticEnv()

	// Setting defaults for this application
	viper.SetDefault("logging.level", "error")

	viper.SetDefault("queue.listenPort", 8888)

	viper.SetDefault("leader.ListenPort", 8080) // Not used as of now
	viper.SetDefault("leader.ReplicaHost", "localhost")
	viper.SetDefault("leader.DataNodeStateCheckPeriod", 30)
	viper.SetDefault("leader.DataNodeRemainingCheckPeriod", 2)

	viper.SetDefault("metric.ListenPort", 9000)

	// Read Config from Flags
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	// Read Config from file
	if configFile, err := cmd.Flags().GetString("config-file"); err == nil && configFile != "" {
		viper.SetConfigFile(configFile)

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func provideConfig(cmd *cobra.Command) (*Config, error) {
	config, err := LoadConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load configurations.")
	}
	return config, nil
}
