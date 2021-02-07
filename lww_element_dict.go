package crdt

type LWWElementDict struct {
	id      int
	clock   Clock
	adds    map[string]TimedValue
	removes map[string]TimedValue
}

func NewLWWElementDict(id int, clock Clock) LWWElementDict {
	return LWWElementDict{
		id:      id,
		clock:   clock,
		adds:    map[string]TimedValue{},
		removes: map[string]TimedValue{},
	}
}

// Add-or-updates a key-value pair, and generates a new timestamp.
// Timestamps are assumed unique, totally ordered, consistent with causal order, and monotonically increasing.
func (l LWWElementDict) Add(key string, value interface{}) {
	timedValue := NewTimedValue(value, l.clock.Now())
	if add, ok := l.adds[key]; ok && add.Compare(timedValue) >= 0 { // consistent with causal order
		return
	}
	l.adds[key] = timedValue
}

// Removes an existing key-value pair, and generates a new timestamp.
// Timestamps are assumed unique, totally ordered, consistent with causal order, and monotonically increasing.
func (l LWWElementDict) Remove(key string) {
	if !l.Lookup(key) { // an element can only be removed if it is present
		return
	}
	timedValue := NewTimedValue(nil, l.clock.Now())
	if remove, ok := l.removes[key]; ok && remove.Compare(timedValue) >= 0 { // consistent with causal order
		return
	}
	l.removes[key] = timedValue
}

// Lookup key-value pair by its key.
// The lookup is biased towards addition.
func (l LWWElementDict) Lookup(key string) bool {
	add, ok := l.adds[key]
	if !ok {
		return false
	}
	remove, ok := l.removes[key]
	if ok && remove.Compare(add) > 0 { // biased towards addition
		return false
	}
	return true
}

// Returns the current value:
// A collection of key-value pairs without their underlying add-and-remove timestamps.
func (l LWWElementDict) Values() map[string]interface{} {
	values := map[string]interface{}{}
	for key, _ := range l.adds {
		if l.Lookup(key) {
			values[key] = l.adds[key].value
		}
	}
	return values
}

// Merge commutes the values with the highest timestamp.
// On concurrent add-or-update, the node with the lowest identifier (address) takes precedence.
// Therefore, concurrent add-or-update operations are deterministic.
// Merge is commutative, associative, and idempotent.
func (l LWWElementDict) Merge(other LWWElementDict) LWWElementDict {
	dominant := l.id >= other.id
	result := NewLWWElementDict(l.id, l.clock)
	result.adds = merge(l.adds, other.adds, dominant)
	result.removes = merge(l.removes, other.removes, dominant)
	return result
}

func merge(source map[string]TimedValue, replica map[string]TimedValue, dominant bool) map[string]TimedValue {
	result := map[string]TimedValue{}
	for key, value := range source {
		result[key] = value
	}
	for key, replicaValue := range replica {
		sourceValue, ok := source[key]
		if !ok {
		} else if sourceValue.Compare(replicaValue) > 0 ||
			(sourceValue.Compare(replicaValue) == 0 && !dominant) {
			continue
		}
		result[key] = replicaValue
	}
	return result
}
