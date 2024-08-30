package placementpool

import (
	"errors"
	"math"
	"math/rand/v2"

	"golang.org/x/text/currency"
)

var ErrEmptyQueue = errors.New("empty queue")

// FISO first in shrug out queue -- items are dequeued when RNGesus sees fit
// uses Fisher-Yates and underlying data is never discarded instead the slice
// is "sorted" in place. Unlike most in place FY implementations, this
// structure moves swapped elements to the front of the slice rather than the
// end. this allows new items to be added after dequeues have happened
func NewRandomQueue[T any](rng *rand.Rand) *FisherYatesQueue[T] {
	r := new(FisherYatesQueue[T])
	r.rng = rng
	return r
}

// the passed slice is immediately shuffled
func RandomQueueFrom[T any](rng *rand.Rand, ts []T) *FisherYatesQueue[T] {
	r := new(FisherYatesQueue[T])
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
type FisherYatesQueue[T any] struct {
	dqCount int
	members []T
	rng     *rand.Rand
}

func (r *FisherYatesQueue[T]) Len() int {
	return len(r.members) - r.dqCount
}

func (r *FisherYatesQueue[T]) Dequeue() (*T, error) {
	curLen := len(r.members)

	if curLen == 0 || r.dqCount == curLen {
		return nil, ErrEmptyQueue
	}

	return r.dequeue(), nil
}

func (r *FisherYatesQueue[T]) dequeue() *T {
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
	return &r.members[current]
}

func (r *FisherYatesQueue[T]) Enqueue(t T) {
	r.members = append(r.members, t)
}

func (r *FisherYatesQueue[T]) EnqueueSlice(ts []T) {
	r.members = append(r.members, ts...)
}

// Requeueing happens when the fill algorithm determines it cannot place the
// provided thing _at this moment_. For example if the fill algorithm dequeues
// "Hover Boots" and the "Shadow Temple Compass Chest" location it will
// determine it cannot place this token and location pair and return both to
// their respective pools. The caller is responsible for knowing how many items
// to requeue but this operation caps the requeues to the first dequeued item.
//
//		```golang
//		q := RandomQueueFrom(...)
//		// safe even if queue is empty
//		q.RequeueLastN(math.IntMax)
//	 ```
//
// However, RequeueAll is provided to both reshuffle the entire queue and
// resets the dqCount count to 0.
func (r *FisherYatesQueue[T]) RequeueLastN(n int) {
	r.dqCount = int(math.Min(0, float64(r.dqCount-n)))
}

// Most common case of requeuing
func (r *FisherYatesQueue[T]) RequeueLast() {
	r.RequeueLastN(1)
}

// Requeues every item in the pool. Unlike calling
// `r.RequeueLastN(math.IntMax)` this will additionally re-shuffle the pool
// to remove any incidental sorting due to previous usage of the queue.
func (r *FisherYatesQueue[T]) RequeueAll() {
	r.dqCount = 0
	r.ShuffleRemaining()
}

// Shuffles only indexes eligble for dequeuing.
func (r *FisherYatesQueue[T]) ShuffleRemaining() {
	r.rng.Shuffle(r.Len(), func(i, j int) {
		r.swap(i+r.dqCount, j+r.dqCount)
	})
}

func (r *FisherYatesQueue[T]) swap(i, j int) {
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
func (r *FisherYatesQueue[T]) All(yield func(*T) bool) {
	for r.Len() > 0 {
		if !yield(r.dequeue()) {
			return
		}
	}
}
