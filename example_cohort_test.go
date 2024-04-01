package chronos_test

import (
	"sync"
	"testing"
	"time"

	"github.com/seborama/chronos"
)

func TestCohort(t *testing.T) {
	coh1 := &chronos.Cohort{}

	var wg sync.WaitGroup

	for i := 1; i <= 10_000; i++ {
		c1 := chronos.Builder{}.WithSamplingRate(100).BuildTS()
		c2 := chronos.Builder{}.WithSamplingRate(100).BuildTS()
		c3 := chronos.Builder{}.WithSamplingRate(100).BuildTS()
		coh1.Add(c1, c2, c3)

		wg.Add(1)
		go func(c1 *chronos.ChronosTS) {
			defer wg.Done()
			for i := 1; i <= 1_000; i++ {
				c1.Start()
				time.Sleep(1 * time.Nanosecond)
				c1.Stop()

				c2.Start()
				time.Sleep(3 * time.Nanosecond)
				c2.Stop()

				c3.Start()
				time.Sleep(2 * time.Nanosecond)
				c3.Stop()
			}
		}(c1)
	}

	wg.Wait()

	coh1.Println("doSomething")
}
