package chronos

import (
	"time"
)

// Chronos is a chronometer.
// This is implementation is NOT thread-safe, see ChronosTS for a thread-safe implementation.
// Chronos offers somewhat improved performance over ChronosTS.
type Chronos struct {
	skip          bool
	elapsed       time.Duration
	clock         time.Time
	clockIsActive bool
	count         uint64
	samplingRate  uint32 // default value must be >= 1
	samplingCount uint32 // default value must be 0
}

func (c *Chronos) Skip() {
	c.skip = true
}

func (c *Chronos) Start() {
	if c.skip {
		return
	}

	c.samplingCount++
	if c.samplingCount != c.samplingRate {
		return
	}

	c.clock = time.Now()
	c.clockIsActive = true
	c.samplingCount = 0
}

func (c *Chronos) Stop() {
	if c.skip {
		return
	}

	if c.samplingCount != 0 {
		if c.clockIsActive {
			panic("Internal Error: clock should not be running when not sampling")
		}
		return
	}

	if !c.clockIsActive {
		panic("cannot stop inactive chronos")
	}

	c.elapsed += time.Since(c.clock)
	c.clockIsActive = false
	c.count++
}

// Elapsed returns the cumulative elapsed time.
// It will NOT wait for the timer to stop.
func (c *Chronos) Elapsed() time.Duration {
	if c.skip {
		return 0
	}
	return c.elapsed
}

// Count returns the count of durations measured.
// It will NOT wait for the timer to stop.
func (c *Chronos) Count() uint64 {
	if c.skip {
		return 0
	}
	return c.count
}

func (c *Chronos) Metrics() Metrics {
	return metrics(c)
}

// Println returns information about this chronometer.
// It will NOT wait for the timer to stop.
func (c *Chronos) Println(label string) {
	if c.skip {
		return
	}

	m := c.Metrics()
	printLn(m, label)
}
