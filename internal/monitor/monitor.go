package monitor

import (
	"log"
	"strings"
	"time"

	"github.com/gaiaz-iusipov/rpi-violet/internal/monitor/co2mon"
)

// Monitor is a USB CO2 monitor
type Monitor struct {
	dev       *co2mon.Device
	readDelay time.Duration
	co2       *co2
	temp      *temp
	done      chan struct{}
}

func New(dev *co2mon.Device, setters ...OptionSetter) *Monitor {
	opts := newOptions(setters)

	monitor := &Monitor{
		dev:       dev,
		readDelay: opts.readDelay,
		temp:      &temp{ttl: opts.co2TTL},
		co2:       &co2{ttl: opts.co2TTL},
		done:      make(chan struct{}),
	}

	go monitor.run()
	return monitor
}

func (m *Monitor) Close() error {
	close(m.done)
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

func (m *Monitor) run() {
	ticker := time.NewTicker(m.readDelay)
	defer ticker.Stop()

	for {
		select {
		case <-m.done:
			return
		case <-ticker.C:
			pack, err := m.dev.ReadPacket()
			if err != nil {
				log.Println(err)
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
