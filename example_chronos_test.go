package chronos_test

import (
	"sync"
	"testing"
	"time"

	"github.com/seborama/chronos"
)

func TestExampleChronos(t *testing.T) {
	ch := chronos.Builder{}.Build()
	ch.Start()
	doSomething()
	ch.Stop()
	ch.Println("doSomething")
}

func doSomething() {
	time.Sleep(time.Millisecond)
}

func TestExampleChronosTS(t *testing.T) {
	coh1 := &chronos.Cohort{}

	var wg sync.WaitGroup

	for i := 1; i <= 10_000; i++ {
		c1 := chronos.Builder{}.WithSamplingRate(100).BuildTS()
		coh1.Add(c1)

		wg.Add(1)
		go func(c1 *chronos.ChronosTS) {
			defer wg.Done()
			for i := 1; i <= 1_000; i++ {
				c1.Start()
				time.Sleep(1 * time.Nanosecond)
				c1.Stop()
			}
		}(c1)
	}

	wg.Wait()

	coh1.Println("doSomething")
}
