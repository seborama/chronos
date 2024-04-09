package chronos

import (
	"fmt"
	"sync"
)

type printer interface {
	Println(label string)
}

// An Outline is an executive summary of the aggregated timing of all the printers
// it holds in order they were added.
// It is typically used to aggregate `Cohort`'s but it can equally aggregate `Chronos`'s.
type Outline struct {
	members []printer
	labels  []string
	lock    sync.Mutex
}

func (c *Outline) Add(labels []string, members []printer) {
	if len(labels) != len(members) {
		panic("unbalanced labels and elements")
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.members = append(c.members, members...)
	c.labels = append(c.labels, labels...)
}

func (c *Outline) Println(label string) {
	if len(c.members) == 0 {
		fmt.Println(label, ">> no data")
		return
	}

	fmt.Println(label, "---")
	for _, member := range c.members {
		fmt.Print("    ")
		member.Println("x")
	}
	fmt.Println("---", label)
}
