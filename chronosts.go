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
	count         int32
	samplingRate  int32
	samplingCount int32
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
func (c *ChronosTS) Count() int32 {
	if c.skip {
		return 0
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.count
}

// Println returns information about this chronometer.
// It will wait for the timer to stop.
func (c *ChronosTS) Println(label string) {
	if c.skip {
		return
	}

	printLn(c, label)
}

func printLn(c chronometer, label string) {
	if c.Count() == 0 {
		fmt.Println(label, ">> no data")
		return
	}

	avg := time.Duration(int64(c.Elapsed()) / int64(c.Count()))
	// c: count
	// e: elapsed
	// avg: average
	fmt.Println(label, ">> c:", c.Count(), "e:", c.Elapsed(), "avg:", avg)
}
