package chronos_test

import (
	"sync"
	"testing"
	"time"

	"github.com/seborama/chronos"
)

func TestOutline(t *testing.T) {
	coh1 := &chronos.Cohort{}
	coh2 := &chronos.Cohort{}

	ol1 := &chronos.Outline{}
	ol1.Add(coh1, coh2)

	var wg sync.WaitGroup

	for i := 1; i <= 1_000; i++ {
		c1 := chronos.Builder{}.WithSamplingRate(100).BuildTS()
		coh1.Add(c1)
		c2 := chronos.Builder{}.WithSamplingRate(10).BuildTS()
		coh2.Add(c2)

		wg.Add(1)
		go func(c1 *chronos.ChronosTS) {
			defer wg.Done()
			for i := 1; i <= 500; i++ {
				c1.Start()
				time.Sleep(1 * time.Nanosecond)
				c1.Stop()

				for j := 1; j <= 100; j++ {
					c2.Start()
					time.Sleep(2 * time.Nanosecond)
					c2.Stop()
				}
			}
		}(c1)
	}

	wg.Wait()

	ol1.Println("testOutline")
}
