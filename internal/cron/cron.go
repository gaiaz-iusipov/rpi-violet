package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"periph.io/x/conn/v3/gpio"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
)

type Cron struct {
	cron *cron.Cron
}

type PhotoProvider interface {
	GetPhoto(ctx context.Context) ([]byte, error)
}

type MeasurementsProvider interface {
	Measurements() (string, bool)
}

type PhotoSender interface {
	SendPhoto(ctx context.Context, photo []byte, caption string) error
}

func New(
	cfg *config.Cron,
	loc *time.Location,
	pin gpio.PinIO,
	photoProvider PhotoProvider,
	photoSender PhotoSender,
	measurementsProvider MeasurementsProvider,
) (*Cron, error) {
	c := cron.New(
		cron.WithLocation(loc),
	)

	for _, jobCfg := range cfg.Jobs {
		_, err := c.AddJob(jobCfg.Spec, newJob(jobCfg, pin, photoProvider, photoSender, measurementsProvider))
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
