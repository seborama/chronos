# chronos

No frills, simple gathering of execution stats. This is essentially no more than a convenient wrapper for `time.Since`.

Has a high footprint on performance but is more sensitive and granular than Go's profiler.

Useful when you have blocks of code that don't span across functions or that execute in just a few nanoseconds.

Since Chronos has a high footprint (the "quantum" conundrum), it is best to avoid nesting Chronos's.

## ChronosTS

For strong thread safety (at a cost on performance), use `ChronosTS`.

## Examples

```go
func main() {
	ch := chronos.Builder{}.Build()
	ch.Start()
	doSomething()
	ch.Stop()
	ch.Println("doSomething")
}

// doSomething >> c: 1 e: 1.1425ms avg: 1.1425ms
```

### Threaded code

```go
func TestExample1(t *testing.T) {
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

// doSomething >> c: 100000 tcd: 260.336807ms avg: 2.603µs mpt: 26.033µs
```
