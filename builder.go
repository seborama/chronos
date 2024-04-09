package chronos

type Builder struct {
	skip         *bool
	samplingRate *uint32
}

func (b Builder) WithSkip() *Builder {
	val := true
	b.skip = &val
	return &b
}

func (b Builder) WithSamplingRate(val uint32) *Builder {
	b.samplingRate = &val
	return &b
}

func (b Builder) Build() *Chronos {
	c := Chronos{}

	if b.skip != nil {
		c.skip = *b.skip
	}

	if b.samplingRate != nil {
		c.samplingRate = *b.samplingRate
	} else {
		c.samplingRate = 1
	}

	return &c
}

func (b Builder) BuildTS() *ChronosTS {
	c := ChronosTS{}

	if b.skip != nil {
		c.skip = *b.skip
	}

	if b.samplingRate != nil {
		c.samplingRate = *b.samplingRate
	} else {
		c.samplingRate = 1
	}

	return &c
}
