package monitor

import (
	"fmt"
	"sync"
	"time"
)

// Temperature in Celsius.
type temp struct {
	value     int
	updatedAt time.Time
	mutex     sync.RWMutex
	ttl       time.Duration
}

func (t *temp) set(value int) {
	now := time.Now()

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.value = value
	t.updatedAt = now
}

func (t *temp) Get() (int, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Check TTL
	if time.Since(t.updatedAt) > t.ttl {
		return 0, false
	}

	return t.value, t.value != 0
}

// getCelsius returns a temperature value in degrees celsius.
func (t *temp) getCelsius() (float64, bool) {
	if val, ok := t.Get(); ok {
		return float64(val)/16.0 - 273.15, true
	}
	return 0, false
}

func (t *temp) String() string {
	if val, ok := t.getCelsius(); ok {
		return fmt.Sprintf("%.4f °C", val)
	}
	return ""
}
