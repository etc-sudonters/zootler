package shuffle

import (
	"errors"
	"math/rand/v2"
)

var ErrEmptyQueue = errors.New("empty queue")
var ErrCannotRequeue = errors.New("cannot requeue item")

// FISO first in shrug out queue -- items are dequeued when RNGesus sees fit
// uses Fisher-Yates and underlying data is never discarded instead the slice
// is "sorted" in place. Unlike most in place FY implementations, this
// structure moves swapped elements to the front of the slice rather than the
// end. this allows new items to be added after dequeues have happened
func Empty[T comparable](rng *rand.Rand) *Q[T] {
	r := new(Q[T])
	r.rng = rng
	return r
}

// the passed slice is immediately shuffled
func From[T comparable](rng *rand.Rand, ts []T) *Q[T] {
	r := new(Q[T])
	r.rng = rng
	r.members = ts
	r.ShuffleRemaining()
	return r
}

// A "queue" that dequeues elements via incremental Fisher-Yates shuffling.
// Unlike most Fisher Yates implementations, this queue swaps items to
// _beginning_ of the slice.  This allows the queue to grow after some items
// have been dequeued. Additionally data is never discarded which provides a
// convenient and fast dequeue and requeue method for the fillers.
type Q[T comparable] struct {
	dqCount int
	members []T
	rng     *rand.Rand
}

func (r *Q[T]) Len() int {
	return len(r.members) - r.dqCount
}

func (r *Q[T]) Dequeue() (T, error) {
	var empty T
	curLen := len(r.members)

	if curLen == 0 || r.dqCount == curLen {
		return empty, ErrEmptyQueue
	}

	return r.dequeue(), nil
}

func (r *Q[T]) dequeue() T {
	// IntN is the half open range 0 (inclusive) to N (exclusive)
	// r.Len() is len(r.members) - r.dqCount and provides N.
	//   note: len(r.members) is the exclusive upper bound, any number 0
	//   (inclusive) to len(..) (exclusive) is a valid index in the slice
	// Since we swap to the _front_ of the slice we add the dqCount count to
	// this generated index to get the actual swapping index.
	// do the swap and spit out a pointer to the selected item
	swap := r.rng.IntN(r.Len()) + r.dqCount
	current := r.dqCount
	r.dqCount += 1
	r.swap(current, swap)
	return r.members[current]
}

func (r *Q[T]) Enqueue(t T) {
	r.members = append(r.members, t)
}

func (r *Q[T]) EnqueueSlice(ts []T) {
	members := make([]T, 0, len(r.members)+len(ts))
	members = append(members, r.members...)
	members = append(members, ts...)
	r.members = members
}

func (r *Q[T]) Requeue(t T) error {
	if r.dqCount == 0 {
		return ErrEmptyQueue
	}

	if r.members[r.dqCount] != t {
		return ErrCannotRequeue
	}

	r.dqCount -= 1
	return nil
}

// Shuffles only indexes eligble for dequeuing.
func (r *Q[T]) ShuffleRemaining() {
	r.rng.Shuffle(r.Len(), func(i, j int) {
		r.swap(i+r.dqCount, j+r.dqCount)
	})
}

func (r *Q[T]) swap(i, j int) {
	r.members[i], r.members[j] = r.members[j], r.members[i]
}

// Convenience iterator that repeatedly randomly dequeues until all items are
// dequeued. Items can be requeued as normal, however care must be taken to
// avoid accidentally creating infinite loops:
//
//	```golang
//	q := RandomQueueFrom(rng, theSlice)
//	for something := range q.All {
//	    if !predicate(something) {
//	        q.RequeueLast()
//	    }
//	}
//	```
//
// Depending on how `predicate` determines something should be requeued it may
// be possible to end up in a situation where the predicate _always_ rejects
// the item, placing it back into the queue and eventually the queue consists
// exclusively of items the predicate will reject.
func (r *Q[T]) All(yield func(T) bool) {
	for r.Len() > 0 {
		if !yield(r.dequeue()) {
			return
		}
	}
}
