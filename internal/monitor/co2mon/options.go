package co2mon

import (
	"crypto/rand"
)

const (
	// defaultPath of a USB device.
	defaultPath = "/dev/hidraw0"
)

type options struct {
	path           string
	key            [8]byte
	withoutDecrypt bool
}

func newOptions(setters []OptionSetter) *options {
	opts := &options{
		path: defaultPath,
	}
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}

type OptionSetter func(opts *options)

// WithPath sets a device file path.
// Default value is `/dev/hidraw0`.
func WithPath(path string) OptionSetter {
	return func(opts *options) {
		opts.path = path
	}
}

// WithKey sets a static key.
func WithKey(key [8]byte) OptionSetter {
	return func(opts *options) {
		opts.key = key
	}
}

// WithRandomKey sets a random key.
func WithRandomKey() OptionSetter {
	return func(opts *options) {
		_, _ = rand.Read(opts.key[:])
	}
}

// WithoutDecrypt disables data decrypt.
func WithoutDecrypt() OptionSetter {
	return func(opts *options) {
		opts.withoutDecrypt = true
	}
}
