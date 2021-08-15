package monitor

import (
	"time"
)

type options struct {
	readDelay       time.Duration
	co2TTL, tempTTL time.Duration
}

func newOptions(setters []OptionSetter) *options {
	opts := &options{
		readDelay: time.Second,
		co2TTL:    10 * time.Second,
		tempTTL:   10 * time.Second,
	}
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}

type OptionSetter func(opts *options)

func WithReadDelay(readDelay time.Duration) OptionSetter {
	return func(opts *options) {
		opts.readDelay = readDelay
	}
}
