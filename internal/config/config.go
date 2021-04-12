package config

import (
	"github.com/go-playground/validator/v10"
)

type Config struct {
	DebugPort uint16 `validate:"required"`
	TimeZone  string `validate:"required"`
	*Sentry
	*GPIO
	*Telegram
	*Raspistill
	*Cron
}

type Sentry struct {
	DNS     string   `validate:"required,url"`
	Timeout Duration `validate:"required"`
}

type GPIO struct {
	LightPin string `validate:"required"`
}

type Telegram struct {
	BotToken      string   `validate:"required"`
	ChatID        int64    `validate:"required"`
	ClientTimeout Duration `validate:"required"`
}

type Raspistill struct {
	Quality uint8    `validate:"min=1,max=100"`
	Timeout Duration `validate:"required"`
}

func (c *Config) validate() error {
	v := validator.New()
	return v.Struct(c)
}

type Cron struct {
	Jobs []*CronJob
}

type CronJob struct {
	Spec            string `validate:"required"`
	WithLightSwitch bool
	LightState      bool
	WithPhoto       bool
}
