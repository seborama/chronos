package chronos

import (
	"fmt"
	"sync"
	"time"
)

type chronometer interface {
	Elapsed() time.Duration
	Count() int32
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

func (c *Cohort) Metrics() (e time.Duration, cnt uint64) {
	for _, ch := range c.members {
		e += ch.Elapsed() // SIGSEGV (????)
		cnt += uint64(ch.Count())
	}

	return
}

func (c *Cohort) Println(label string) {
	e, cnt := c.Metrics()

	if cnt == 0 {
		fmt.Println(label, ">> no data")
		return
	}

	avg := time.Duration(int64(e) / int64(cnt))
	// tcd: total cumulative duration
	// mpt: mean per thread
	mpt := time.Duration(int64(e) / int64(len(c.members)))
	fmt.Println(label, ">> c:", cnt, "tcd:", e, "avg:", avg, "mpt:", mpt)
}
