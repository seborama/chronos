package chronos

import (
	"fmt"
	"sync"
	"time"
)

type chronometer interface {
	Elapsed() time.Duration
	Count() uint64
	Println(label string)
}

// A Cohort represents a group of chronometers that should behave as one.
// A typical usage is an application that generates multiple threads of the same
// go routine.
// Just like a Chronos reports an **aggregation** of all the measures it performed,
// a Cohort reports an aggregation of the all the chronometers it holds.
type Cohort struct {
	members []chronometer
	lock    sync.Mutex
}

func (c *Cohort) Add(els ...chronometer) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.members = append(c.members, els...)
}

type Metrics struct {
	Elapsed time.Duration
	Count   uint64
	Average time.Duration
}

func (c *Cohort) Metrics() Metrics {
	var (
		e   time.Duration
		cnt uint64
	)

	for _, ch := range c.members {
		e += ch.Elapsed() // SIGSEGV (????)
		cnt += uint64(ch.Count())
	}

	var avg time.Duration
	if cnt > 0 {
		avg = time.Duration(int64(e) / int64(cnt))
	}

	return Metrics{
		Elapsed: e,
		Count:   cnt,
		Average: avg,
	}
}

func (c *Cohort) Println(label string) {
	m := c.Metrics()

	if m.Count == 0 {
		fmt.Println(label, ">> no data")
		return
	}

	meanPerThread := time.Duration(int64(m.Elapsed) / int64(len(c.members)))
	fmt.Println(label, ">> count:", m.Count, "total-cumulative-duration:", m.Elapsed, "avg:", m.Average, "mean-per-thread:", meanPerThread)
}
