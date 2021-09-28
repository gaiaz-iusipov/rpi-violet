package monitor

import (
	"time"

	"github.com/gaiaz-iusipov/rpi-violet/internal/monitor/co2mon"
)

type options struct {
	readDelay        time.Duration
	co2TTL, tempTTL  time.Duration
	readErrThreshold int
	devOpts          []co2mon.OptionSetter
}

func newOptions(setters []OptionSetter) *options {
	opts := &options{
		readDelay:        time.Second,
		co2TTL:           10 * time.Second,
		tempTTL:          10 * time.Second,
		readErrThreshold: 5,
	}
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}

type OptionSetter func(opts *options)

func WithDevOptions(setters ...co2mon.OptionSetter) OptionSetter {
	return func(opts *options) {
		opts.devOpts = setters
	}
}

func WithReadDelay(readDelay time.Duration) OptionSetter {
	return func(opts *options) {
		opts.readDelay = readDelay
	}
}
