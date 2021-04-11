package config

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	TelegramBotToken      string
	TelegramClientTimeout time.Duration
	TelegramChatID        int64
	GPIOLightPin          string
	RaspistillTimeout     time.Duration
	SentryDNS             string
	SentryTimeout         time.Duration
	DebugPort             string
}

func (c *Config) validate() error {
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
