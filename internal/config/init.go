package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func Init(cfgFile string) (*Config, error) {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	log.Println("using config file:", viper.ConfigFileUsed())
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var cfg *Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}
