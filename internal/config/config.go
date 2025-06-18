package config

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	ListenAddress string   `mapstructure:"listen_address"`
	Database      Database `mapstructure:"db"`
}

func New(ctx context.Context, cfgFile string) *Config {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("hms-config") // name of config file (without extension)
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME") // adding home directory as first search path
		viper.AddConfigPath(".")     // optionally look for config in the working directory
	}

	viper.SetDefault("listen_address", ":8080")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.name", "hatch_messaging_service")

	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.WarnContext(ctx, fmt.Sprintf("Config file not found; ignore error if desired: %v", err))
		} else {
			slog.ErrorContext(ctx, fmt.Sprintf("Config file was found but another error was produced: %v", err))
		}
	}

	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("unable to unmarshal config into struct, %v", err))
	}

	return &config
}
