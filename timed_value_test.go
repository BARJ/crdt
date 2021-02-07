package crdt

import "testing"

func TestTimedValueCompare(t *testing.T) {
	tests := []struct {
		name        string
		timedValue1 TimedValue
		timedValue2 TimedValue
		wantCompare int
	}{
		{
			name:        "TimedValue.timestamp == TimedValue.timestamp",
			timedValue1: NewTimedValue("foo", 10),
			timedValue2: NewTimedValue("bar", 10),
			wantCompare: 0,
		},
		{
			name:        "TimedValue.timestamp > TimedValue.timestamp",
			timedValue1: NewTimedValue("foo", 15),
			timedValue2: NewTimedValue("bar", 10),
			wantCompare: 1,
		},
		{
			name:        "TimedValue.timestamp < TimedValue.timestamp",
			timedValue1: NewTimedValue("foo", 10),
			timedValue2: NewTimedValue("bar", 15),
			wantCompare: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotCompare := test.timedValue1.Compare(test.timedValue2)
			if gotCompare != test.wantCompare {
				t.Errorf("got %d, want %d", gotCompare, test.wantCompare)
			}
		})
	}
}
