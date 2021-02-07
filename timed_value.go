package crdt

import "fmt"

type TimedValue struct {
	value     interface{}
	timestamp int64
}

func NewTimedValue(value interface{}, timestamp int64) TimedValue {
	return TimedValue{value: value, timestamp: timestamp}
}

func (v TimedValue) Compare(other TimedValue) int {
	if v.timestamp == other.timestamp {
		return 0
	}
	if v.timestamp < other.timestamp {
		return -1
	}
	return +1
}

func (v TimedValue) String() string {
	return fmt.Sprintf("{value=%v, timestamp=%d}", v.value, v.timestamp)
}
