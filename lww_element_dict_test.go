package crdt

import (
	"testing"
)

func TestLWWElementDictAdd(t *testing.T) {
	clock := NewFakeClock()
	l := NewLWWElementDict(1, clock)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})

	// add new key-value pair to an empty state
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})
	clock.SetTime(1)
	l.Add("a", 1)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})

	// add new key-value pair to an non-empty state
	l.Add("b", "foo")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1, "b": "foo"})

	// update existing key-value pair
	clock.SetTime(2)
	l.Add("b", "bar")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1, "b": "bar"})

	// update existing key-value pair with same timestamp: no-op
	l.Add("b", "zi")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1, "b": "bar"})

	// update existing key-value pair with lesser timestamp: no-op
	clock.SetTime(1)
	l.Add("b", "zi")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1, "b": "bar"})

	// re-add key-value pair
	clock.SetTime(3)
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"b": "bar"})
	clock.SetTime(4)
	l.Add("a", 1)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1, "b": "bar"})
}

func TestLWWElementDictRemove(t *testing.T) {
	clock := NewFakeClock()
	l := NewLWWElementDict(1, clock)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})

	// remove non-existing key-value pair
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})
	clock.SetTime(1)
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})

	// remove existing key-value pair
	clock.SetTime(2)
	l.Add("a", 1)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})
	clock.SetTime(3)
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})

	// re-remove existing key-value pair
	clock.SetTime(4)
	l.Add("a", 1)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})
	clock.SetTime(5)
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})

	// re-remove existing key-value pair with same timestamp: no-op
	clock.SetTime(6)
	l.Add("a", 1)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})

	// re-remove existing key-value pair with lesser timestamp: no-op
	clock.SetTime(5)
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})
}

func TestLWWElementDictLookup(t *testing.T) {
	clock := NewFakeClock()
	l := NewLWWElementDict(1, clock)

	// lookup non-existing key-value pair
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})
	if l.Lookup("a") {
		t.Error("unexpected key-value pair")
	}

	// lookup existing key-value pair
	clock.SetTime(1)
	l.Add("a", 1)
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})
	if !l.Lookup("a") {
		t.Error("expected key-value pair")
	}

	// lookup removed key-value pair
	clock.SetTime(2)
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{})
	if l.Lookup("a") {
		t.Error("unexpected key-value pair")
	}

	// lookup concurrent add/remove key-value pair: biased towards addition
	clock.SetTime(3)
	l.Add("a", 1)
	l.Remove("a")
	assertLWWElementDict(t, l.Values(), map[string]interface{}{"a": 1})
	if !l.Lookup("a") {
		t.Error("expected key-value pair")
	}
}

func TestLWWElementDictMerge(t *testing.T) {
	clock := NewFakeClock()

	// commutative: l1 U l2 = l2 U l1
	l1 := NewLWWElementDict(1, clock)
	l2 := NewLWWElementDict(2, clock)
	l1.Add("a", "foo")
	l2.Add("b", "bar")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"b": "bar"})
	m1 := l1.Merge(l2)
	m2 := l2.Merge(l1)
	assertLWWElementDict(t, m1.Values(), m2.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{"a": "foo", "b": "bar"})

	// idempotent: l1 U l1 = l1
	l1 = NewLWWElementDict(1, clock)
	l1.Add("a", "foo")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	m1 = l1.Merge(l1)
	assertLWWElementDict(t, m1.Values(), l1.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{"a": "foo"})

	// associative: (l1 U l2) U l3 = l1 U (l2 U l3)
	l1 = NewLWWElementDict(1, clock)
	l2 = NewLWWElementDict(2, clock)
	l3 := NewLWWElementDict(3, clock)
	l1.Add("a", "foo")
	l2.Add("b", "bar")
	l3.Add("c", "zi")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"b": "bar"})
	assertLWWElementDict(t, l3.Values(), map[string]interface{}{"c": "zi"})
	m1 = l1.Merge(l2).Merge(l3)
	m2 = l1.Merge(l2.Merge(l3))
	assertLWWElementDict(t, m1.Values(), m2.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{"a": "foo", "b": "bar", "c": "zi"})

	// add commutes
	l1 = NewLWWElementDict(1, clock)
	l1.Add("a", "foo")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	l2 = NewLWWElementDict(2, clock).Merge(l1)
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "foo"})

	// remove commutes
	l1 = NewLWWElementDict(1, clock)
	l1.Add("a", "foo")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	l2 = NewLWWElementDict(2, clock).Merge(l1)
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "foo"})
	clock.SetTime(2)
	l2.Remove("a")
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{})
	l1 = l1.Merge(l2)
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{})

	// concurrent add (or update) key-value pair: biased towards node with the lowest identifier
	l1 = NewLWWElementDict(1, clock)
	l2 = NewLWWElementDict(2, clock)
	l1.Add("a", "foo")
	l2.Add("a", "bar")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "bar"})
	// Order of nodes, source vs replica, does not matter.
	// The node with the lowest identifier always wins,
	// therefore, concurrent add (or update) operations are deterministic
	m1 = l1.Merge(l2)
	m2 = l2.Merge(l1)
	assertLWWElementDict(t, m1.Values(), m2.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{"a": "foo"})

	// last add wins
	l1 = NewLWWElementDict(1, clock)
	l2 = NewLWWElementDict(2, clock)
	clock.SetTime(1)
	l1.Add("a", "foo")
	clock.SetTime(2)
	l2.Add("a", "bar")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "bar"})
	m1 = l1.Merge(l2)
	m2 = l2.Merge(l1)
	assertLWWElementDict(t, m1.Values(), m2.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{"a": "bar"})

	// late add wins (add.t > remove.t)
	l1 = NewLWWElementDict(1, clock)
	clock.SetTime(1)
	l1.Add("a", "foo")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	l2 = NewLWWElementDict(2, clock).Merge(l1)
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "foo"})
	clock.SetTime(5)
	l1.Remove("a")
	clock.SetTime(8)
	l2.Add("a", "bar")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{})
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "bar"})
	m1 = l1.Merge(l2)
	m2 = l2.Merge(l1)
	assertLWWElementDict(t, m1.Values(), m2.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{"a": "bar"})

	// concurrent add and remove, add wins (add.t == remove.t)
	l1 = NewLWWElementDict(1, clock)
	clock.SetTime(1)
	l1.Add("a", "foo")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	l2 = NewLWWElementDict(2, clock).Merge(l1)
	l2.Remove("a")
	// locally already biased towards addition
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "foo"})
	m1 = l1.Merge(l2)
	m2 = l2.Merge(l1)
	assertLWWElementDict(t, m1.Values(), m2.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{"a": "foo"})

	// late remove wins (remove.t > add.t)
	l1 = NewLWWElementDict(1, clock)
	clock.SetTime(1)
	l1.Add("a", "foo")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{"a": "foo"})
	l2 = NewLWWElementDict(2, clock).Merge(l1)
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "foo"})
	clock.SetTime(8)
	l1.Remove("a")
	clock.SetTime(5)
	l2.Add("a", "bar")
	assertLWWElementDict(t, l1.Values(), map[string]interface{}{})
	assertLWWElementDict(t, l2.Values(), map[string]interface{}{"a": "bar"})
	m1 = l1.Merge(l2)
	m2 = l2.Merge(l1)
	assertLWWElementDict(t, m1.Values(), m2.Values())
	assertLWWElementDict(t, m1.Values(), map[string]interface{}{})
}

func assertLWWElementDict(t *testing.T, gotValues, wantValues map[string]interface{}) {
	t.Helper()

	if len(gotValues) != len(wantValues) {
		t.Errorf("wrong number of elements: got %d, want %d", len(gotValues), len(wantValues))
	}

	for key, gotValue := range gotValues {
		if wantValue, ok := wantValues[key]; !ok {
			t.Errorf("unexpected key-value pair: got %s=%v, want nothing", key, gotValue)
		} else if gotValue != wantValue {
			t.Errorf("wrong value: got %s=%v, want %s=%v", key, gotValue, key, wantValue)
		}
	}

	for key, wantValue := range wantValues {
		if _, ok := gotValues[key]; !ok {
			t.Errorf("missing key-value pair: got nothing, want %s=%v", key, wantValue)
		}
	}
}
