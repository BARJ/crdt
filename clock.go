package crdt

type Clock interface {
	Now() int64
}

var _ Clock = (*FakeClock)(nil)

type FakeClock struct {
	time int64
}

func NewFakeClock() *FakeClock {
	return &FakeClock{}
}

func (c *FakeClock) SetTime(time int64) {
	c.time = time
}

func (c *FakeClock) Now() int64 {
	return c.time
}
