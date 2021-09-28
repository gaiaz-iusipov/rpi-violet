package monitor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gaiaz-iusipov/rpi-violet/internal/monitor/co2mon"
)

// Monitor is a USB CO2 monitor
type Monitor struct {
	opts         *options
	co2          *co2
	temp         *temp
	done         chan struct{}
	readErrCount int
	dev          *co2mon.Device
}

func New(setters ...OptionSetter) (*Monitor, error) {
	opts := newOptions(setters)

	monitor := &Monitor{
		opts: opts,
		temp: &temp{ttl: opts.co2TTL},
		co2:  &co2{ttl: opts.co2TTL},
		done: make(chan struct{}),
	}

	if err := monitor.initDev(); err != nil {
		return nil, err
	}

	go monitor.run()

	return monitor, nil
}

func (m *Monitor) Close() error {
	close(m.done)
	_ = m.dev.Close()
	return nil
}

func (m *Monitor) Measurements() (string, bool) {
	parts := make([]string, 0, 2)

	if val := m.temp.String(); val != "" {
		parts = append(parts, val)
	}

	if val := m.co2.String(); val != "" {
		parts = append(parts, val)
	}

	if len(parts) == 0 {
		return "", false
	}
	return strings.Join(parts, ", "), true
}

func (m *Monitor) initDev() error {
	dev, err := co2mon.Open(m.opts.devOpts...)
	if err != nil {
		return fmt.Errorf("co2mon.Open(): %w", err)
	}
	m.dev = dev
	return nil
}

func (m *Monitor) run() {
	ticker := time.NewTicker(m.opts.readDelay)
	defer ticker.Stop()

	for {
		select {
		case <-m.done:
			return
		case <-ticker.C:
			pack, err := m.dev.ReadPacket()
			if err != nil {
				log.Println(err)
				m.readErrCount++

				if m.readErrCount > m.opts.readErrThreshold {
					err = m.dev.Close()
					if err != nil {
						log.Println(err)
					}

					err = m.initDev()
					if err != nil {
						log.Println(err)
					} else {
						m.readErrCount = 0
					}
				}

				continue
			}

			switch pack.Ops {
			case co2mon.Co2Ops:
				m.co2.set(pack.Value)
			case co2mon.TempOps:
				m.temp.set(pack.Value)
			}
		}
	}
}
