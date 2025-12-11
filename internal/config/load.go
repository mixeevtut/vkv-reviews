package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Load() (App, error) {
	var cfg App

	configPath := pflag.String("config", "config.yaml", "path to config file")

	viper.SetConfigFile(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return cfg, fmt.Errorf("bind flags: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal: %w", err)
	}

	return cfg, nil
}
