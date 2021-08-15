package monitor

import (
	"fmt"
	"sync"
	"time"
)

// co2 carbon dioxide level in ppm.
type co2 struct {
	value     int
	updatedAt time.Time
	mutex     sync.RWMutex
	ttl       time.Duration
}

func (c *co2) set(value int) {
	now := time.Now()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.value = value
	c.updatedAt = now
}

func (c *co2) get() (int, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Check TTL
	if time.Since(c.updatedAt) > c.ttl {
		return 0, false
	}

	return c.value, c.value != 0
}

func (c *co2) String() string {
	if val, ok := c.get(); ok {
		return fmt.Sprintf("%d ppm CO₂", val)
	}
	return ""
}
