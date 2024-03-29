package config

import (
	"fmt"
	"os"

	"github.com/naoina/toml"
)

func New(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}
	defer f.Close()

	config := new(Config)
	err = toml.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	err = config.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return config, nil
}
