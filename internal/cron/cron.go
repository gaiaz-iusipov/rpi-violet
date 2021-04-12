package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"periph.io/x/periph/conn/gpio"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
)

type Cron struct {
	cron *cron.Cron
}

type PhotoProvider interface {
	GetPhoto(context.Context) ([]byte, error)
}

type PhotoSender interface {
	SendPhoto(context.Context, []byte) error
}

func New(cfg *config.Cron, loc *time.Location, pin gpio.PinIO, photoProvider PhotoProvider, photoSender PhotoSender) (*Cron, error) {
	c := cron.New(
		cron.WithLocation(loc),
	)

	for _, jobCfg := range cfg.Jobs {
		_, err := c.AddJob(newJob(jobCfg, pin, photoProvider, photoSender))
		if err != nil {
			return nil, fmt.Errorf("c.AddJob: %w", err)
		}
	}

	c.Start()
	return &Cron{
		cron: c,
	}, nil
}

func (c *Cron) Stop() {
	c.cron.Stop()
}
