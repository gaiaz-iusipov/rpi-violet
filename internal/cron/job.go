package cron

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"periph.io/x/periph/conn/gpio"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
)

type job struct {
	cfg           *config.CronJob
	pin           gpio.PinIO
	photoProvider PhotoProvider
	photoSender   PhotoSender
}

func newJob(cfg *config.CronJob, pin gpio.PinIO, photoProvider PhotoProvider, photoSender PhotoSender) (string, *job) {
	return cfg.Spec, &job{
		cfg:           cfg,
		pin:           pin,
		photoProvider: photoProvider,
		photoSender:   photoSender,
	}
}

func (cj *job) Run() {
	if err := cj.runE(); err != nil {
		localHub := sentry.CurrentHub().Clone()
		localHub.CaptureException(err)
	}
}

func (cj *job) runE() error {
	if cj.cfg.WithLightSwitch {
		pinLvl := gpio.Level(cj.cfg.LightState)
		if err := cj.pin.Out(pinLvl); err != nil {
			return fmt.Errorf("pin.Out: %w", err)
		}
	}

	if cj.cfg.WithPhoto {
		ctx := context.Background()

		photo, err := cj.photoProvider.GetPhoto(ctx)
		if err != nil {
			return fmt.Errorf("failed to get photo: %w", err)
		}

		err = cj.photoSender.SendPhoto(ctx, photo)
		if err != nil {
			return fmt.Errorf("failed to send photo: %w", err)
		}
	}
	return nil
}
