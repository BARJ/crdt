# Conflict-free replicated data type (CRDT)

## State-based Last-Writer-Wins-Element-Dictionary (LWW-Element-Dict)

LWW-Element-Dict is a state-based conflict-free replicated data type with an internal "add-dict" and "remove-dict".
Elements (key-value pairs) are added to the LWW-Element-Dict by adding them to the add-dict. Elements can be removed from the LWW-Element-Dict by adding them to the remove-dict. Only existing elements can be removed from the LWW-Element-Dict. LWW-Element-Dict contains an element if the element is within the add-dict and not within the remove-dict with a timestamp greater than the one in the add-dict. Updating an element is similar to adding an element to the LWW-Element-Dict. When timestamps are equal, add-dict and remove-dict contain the same element with the same timestamp, addition takes precedence, i.e. this implementation is biased towards addition.

Merging one or more replicas of the LWW-Element-Dict consists of taking the union of the add-dicts and the union of the remove-dicts. Merge is commutative, associative, idempotent, and deterministic.

Timestamps are assumed unique, totally ordered, consistent with causal order, and monotonically increasing.

On a local state, before merging, an update (add, update, or remove an element) executes immediately, unaffected by network latency, fault, or disconnection. An update happens entirely at the source and eventually propgates, via the merge operation, to other replicas.

### Operations

Add:
```golang
// Add-or-updates a key-value pair, and generates a new timestamp.
// Timestamps are assumed unique, totally ordered, consistent with causal order, and monotonically increasing.
func (l LWWElementDict) Add(key string, value interface{})
```

Remove:
```golang
// Removes an existing key-value pair, and generates a new timestamp.
// Timestamps are assumed unique, totally ordered, consistent with causal order, and monotonically increasing.
func (l LWWElementDict) Remove(key string)
```

Lookup:
```golang
// Lookup key-value pair by its key.
// The lookup is biased towards addition.
func (l LWWElementDict) Lookup(key string) bool
```

Values:
```golang
// Returns the current value:
// A collection of key-value pairs without their underlying add-and-remove timestamps.
func (l LWWElementDict) Values() map[string]interface{}
```

Merge:
```golang
// Merge commutes the values with the highest timestamp.
// On concurrent add-or-update, the node with the lowest identifier (address) takes precedence.
// Therefore, concurrent add-or-update operations are deterministic.
// Merge is commutative, associative, and idempotent.
func (l LWWElementDict) Merge(other LWWElementDict) LWWElementDict
```

### LWW-Element-Dict vs LWW-Element-Set Challenges

#### Concurrent Add Operations

With set, concurrent add operations are deterministic and do converge:

```
node1.Add('a', t1) U node2.Add('a', t1)  => {'a'}
```

With dict, concurrent add (or update) operations are not deterministic and do not converge:

```
node1.Add('a', 'foo', t1) U node2.Add('a', 'bar', t1) => {'a': ?}
```

Solution: Either have one node take precedence over the other or retain both values.


## Test

### Coverage

See `coverage.html` (97.9%).

### Usage

Run on your local machine (go version go1.15):

```
go test ./... -cover
```

Dockerized:

```
docker build -t barj-crdt .
```
