package commands

import (
	"fmt"
	"log"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	TelegramBotToken      string
	TelegramClientTimeout time.Duration
	TelegramChatID        int64
	GPIOLightPin          uint8
	RaspistillTimeout     time.Duration
	SentryDNS             string
	SentryTimeout         time.Duration
	DebugPort             string
}

func (c *config) validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.TelegramBotToken, validation.Required),
		validation.Field(&c.TelegramClientTimeout, validation.Required),
		validation.Field(&c.TelegramChatID, validation.Required),
		validation.Field(&c.GPIOLightPin, validation.Required),
		validation.Field(&c.RaspistillTimeout, validation.Required),
		validation.Field(&c.SentryDNS, validation.Required, is.URL),
		validation.Field(&c.SentryTimeout, validation.Required),
		validation.Field(&c.DebugPort, validation.Required),
	)
}

var (
	cfgFile string
	cfg     config
)

func initConfig(_ *cobra.Command, _ []string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	log.Println("using config file:", viper.ConfigFileUsed())
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("unable to read config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return nil
}
