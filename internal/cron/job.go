package cron

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"periph.io/x/conn/v3/gpio"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
)

type job struct {
	sentryHub            *sentry.Hub
	cfg                  *config.CronJob
	pin                  gpio.PinIO
	photoProvider        PhotoProvider
	photoSender          PhotoSender
	measurementsProvider MeasurementsProvider
}

func newJob(
	cfg *config.CronJob,
	pin gpio.PinIO,
	photoProvider PhotoProvider,
	photoSender PhotoSender,
	measurementsProvider MeasurementsProvider,
) *job {
	sentryHub := sentry.CurrentHub().Clone()
	sentryHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("cronJobSpec", cfg.Spec)
	})

	return &job{
		sentryHub:            sentryHub,
		cfg:                  cfg,
		pin:                  pin,
		photoProvider:        photoProvider,
		photoSender:          photoSender,
		measurementsProvider: measurementsProvider,
	}
}

func (j *job) Run() {
	if err := j.runE(); err != nil {
		j.sentryHub.CaptureException(err)
	}
}

func (j *job) runE() error {
	if j.cfg.WithLightSwitch {
		pinLvl := gpio.Level(j.cfg.LightState)
		if err := j.pin.Out(pinLvl); err != nil {
			return fmt.Errorf("pin.Out: %w", err)
		}
	}

	if j.cfg.WithPhoto {
		ctx := context.Background()

		photo, err := j.photoProvider.GetPhoto(ctx)
		if err != nil {
			return fmt.Errorf("failed to get photo: %w", err)
		}

		caption, _ := j.measurementsProvider.Measurements()

		err = j.photoSender.SendPhoto(ctx, photo, caption)
		if err != nil {
			return fmt.Errorf("failed to send photo: %w", err)
		}
	}
	return nil
}
