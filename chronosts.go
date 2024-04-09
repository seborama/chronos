package chronos

import (
	"fmt"
	"sync"
	"time"
)

// ChronosTS is the thread-safe implementation of Chronos.
type ChronosTS struct {
	lock          sync.Mutex
	skip          bool
	elapsed       time.Duration
	clock         time.Time
	clockIsActive bool
	count         uint64
	samplingRate  uint32 // default value must be >= 1
	samplingCount uint32 // default value must be 0
}

func (c *ChronosTS) Skip() {
	c.skip = true
}

func (c *ChronosTS) Start() {
	if c.skip {
		return
	}

	c.lock.Lock()

	c.samplingCount++
	if c.samplingCount != c.samplingRate {
		return
	}

	c.clock = time.Now()
	c.clockIsActive = true
	c.samplingCount = 0
}

func (c *ChronosTS) Stop() {
	if c.skip {
		return
	}

	if c.samplingCount != 0 {
		if c.clockIsActive {
			panic("Internal Error: clock should not be active when not sampling")
		}
		c.lock.Unlock()
		return
	}

	if !c.clockIsActive {
		panic("cannot stop inactive chronos")
	}

	c.elapsed += time.Since(c.clock)
	c.clockIsActive = false
	c.count++
	c.lock.Unlock()
}

// Elapsed returns the cumulative elapsed time.
// It will wait for the timer to stop.
func (c *ChronosTS) Elapsed() time.Duration {
	if c.skip {
		return 0
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.elapsed
}

// Count returns the count of durations measured.
// It will wait for the timer to stop.
func (c *ChronosTS) Count() uint64 {
	if c.skip {
		return 0
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.count
}

func (c *ChronosTS) Metrics() Metrics {
	return metrics(c)
}

// Println returns information about this chronometer.
// It will wait for the timer to stop.
func (c *ChronosTS) Println(label string) {
	if c.skip {
		return
	}

	m := c.Metrics()
	printLn(m, label)
}

func metrics(c chronometer) Metrics {
	e := c.Elapsed()
	cnt := c.Count()
	avg := time.Duration(int64(e) / int64(cnt))

	return Metrics{
		Elapsed: e,
		Count:   cnt,
		Average: avg,
	}
}

func printLn(m Metrics, label string) {
	if m.Count == 0 {
		fmt.Println(label, ">> no data")
		return
	}

	fmt.Println(label, ">> count:", m.Count, "elapsed:", m.Elapsed, "avg:", m.Average)
}
