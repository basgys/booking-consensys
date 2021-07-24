package set

import (
	"github.com/basgys/booking-consensys/pkg/timeutil/interval"
)

// Set is an interface implemented by DisjointSet and ImmutableSet.
// It is used when one of these types type must take a set as an argument.
type Set interface {
	// All  returns an ordered slice of all the intervals in the set.
	All() []interval.UnaryInterval
	// Between iterates over the intervals within extents set and calls f
	// with each. If f returns false, iteration ceases.
	//
	// Any interval within the set that overlaps partially with extents is truncated
	// before being passed to f.
	Between(extents interval.UnaryInterval) IntervalRange

	// Contains reports whether x is entirely contained by the set. This means there
	// is no gap.
	//
	// e.g.
	//  Input: [===================]
	//    Set: [====]   |   |      |
	//         |    [===]   |      |
	//         |        |   [===]  |
	//         |        |   |   [==]
	//         |        |gap|      |
	//
	// Output: false
	Contains(x interval.UnaryInterval) bool
	// Overlap returns whether the interval overlaps with other intervals from the Set
	// Note: The relations End Touching and Start Touching are not classified as
	// an overlap
	Overlap(ival interval.UnaryInterval) bool

	// Extent returns the Interval defined by the minimum and maximum values of
	// the set.
	Extent() interval.UnaryInterval
	// Size returns number of intervals in this set
	Size() int
}

type MutableSet interface {
	Set

	// Add adds all the elements of another set to this set.
	Add(b Set)
	// Sub destructively modifies the set by subtracting b.
	Sub(b Set)
	// Intersect destructively modifies the set by intersecting it with b.
	Intersect(b Set)

	// Insert adds interval b to this set
	Insert(b interval.UnaryInterval)
	// Remove removes interval b from this set
	Remove(b interval.UnaryInterval)
}

type IntervalRange interface {
	GetSlice() []interval.UnaryInterval
	Iterator() IntervalIterator
}

type IntervalIterator interface {
	Advance() bool
	Get() interval.UnaryInterval
}

// IntervalReceiver is a function used for iterating over a set of intervals. It
// takes the start and end times and returns true if the iteration should
// continue.
type IntervalReceiver func(interval.UnaryInterval) bool

// Intersect calls fn for each intersecting interval between sets a and b
func Intersect(a, b Set, fn IntervalReceiver) {
	iterX := a.Between(b.Extent()).Iterator()
	iterY := b.Between(a.Extent()).Iterator()

	// Loop through corresponding intervals of S and B.
	// If y == nil, all of the remaining intervals in S are to the right of B.
	// If x == nil, all of the remaining intervals in B are to the right of S.
	moreX := iterX.Advance()
	moreY := iterY.Advance()
	x := iterX.Get()
	y := iterY.Get()
	for moreX && moreY {
		if x.Before(y) {
			moreX = iterX.Advance()
			if moreX {
				x = iterX.Get()
			}
			continue
		}
		if y.Before(x) {
			moreY = iterY.Advance()
			if moreY {
				y = iterY.Get()
			}
			continue
		}
		ival := x.Intersect(y)
		if ival.IsZero() {
			continue
		}

		if !fn(ival) {
			return
		}

		_, right := x.Bisect(y)
		if !right.IsZero() {
			x = right
		} else {
			moreX = iterX.Advance()
			if moreX {
				x = iterX.Get()
			}
		}
	}
}
