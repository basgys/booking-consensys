// Source: A google project
package disjoint

import (
	"fmt"
	"sort"
	"strings"

	"github.com/basgys/booking-consensys/pkg/timeutil/interval"
	"github.com/basgys/booking-consensys/pkg/timeutil/set"
)

// Set is a disjoint set of `interval.UnaryInterval`
type Set struct {
	// non-overlapping intervals
	intervals []interval.UnaryInterval
}

// NewSet returns a new set given a sorted slice of intervals. This function
// panics if the intervals are not sorted.
func NewSet(intervals ...interval.UnaryInterval) *Set {
	for i := 0; i < len(intervals)-1; i++ {
		if !intervals[i].Before(intervals[i+1]) {
			panic(fmt.Errorf(
				"!intervals[%d].Before(intervals[%d]) for %s, %s",
				i, i+1, intervals[i], intervals[i+1],
			))
		}
	}
	return &Set{intervals}
}

// Empty returns a new, empty set of intervals.
func Empty() *Set {
	return &Set{}
}

// Copy returns a copy of a set that may be mutated without affecting the original.
func (s *Set) Copy() *Set {
	return &Set{append([]interval.UnaryInterval(nil), s.intervals...)}
}

// Extent returns the Interval defined by the minimum and maximum values of the
// set.
func (s *Set) Extent() interval.UnaryInterval {
	if len(s.intervals) == 0 {
		return nil
	}
	return s.intervals[0].Encompass(s.intervals[len(s.intervals)-1])
}

// Add adds all the elements of another set to this set.
func (s *Set) Add(b set.Set) {
	it := b.Between(b.Extent()).Iterator()
	for it.Advance() {
		s.Insert(it.Get())
	}
}

func (s *Set) Insert(iv interval.UnaryInterval) {
	if s.Contains(iv) {
		return
	}
	// TODO: Something like Java's ArrayList would allow both O(log(n))
	// insertion and O(log(n)) lookup. For now, we have O(log(n)) lookup and O(n)
	// insertion.
	var newIntervals []interval.UnaryInterval
	push := func(x interval.UnaryInterval) {
		newIntervals = adjoinOrAppend(newIntervals, x)
	}
	inserted := false
	for _, x := range s.intervals {
		if inserted {
			push(x)
			continue
		}
		if iv.Before(x) {
			push(iv)
			push(x)
			inserted = true
			continue
		}
		// [===left===)[==x===)[===right===)
		left, right := iv.Bisect(x)
		if !left.IsZero() {
			push(left)
		}
		push(x)
		// Replace the interval being inserted with the remaining portion of the
		// interval to be inserted.
		if right.IsZero() {
			inserted = true
		} else {
			iv = right
		}
	}
	if !inserted {
		push(iv)
	}
	s.intervals = newIntervals
}

func (s *Set) Remove(iv interval.UnaryInterval) {
	// TODO: Improve this
	s.Sub(NewSet(iv))
}

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
func (s *Set) Contains(x interval.UnaryInterval) bool {
	next := s.iterator(x, true)
	for ival := next(); ival != nil; ival = next() {
		left, right := x.Bisect(ival)
		if !left.IsZero() {
			return false
		}
		x = right
	}
	return x.IsZero()
}

func (s *Set) Overlap(ival interval.UnaryInterval) bool {
	n := sort.Search(len(s.intervals), func(i int) bool {
		return s.intervals[i].Starting().Cmp(ival.Starting()) >= 0
	})
	return s.intervals[min(n, len(s.intervals)-1)].Overlap(ival)
}

// Sub destructively modifies the set by subtracting b.
func (s *Set) Sub(b set.Set) {
	if len(s.intervals) == 0 {
		return
	}

	var intervals []interval.UnaryInterval
	ix := s.Between(s.Extent()).Iterator()
	iy := b.Between(s.Extent()).Iterator()

	ix.Advance()
	iy.Advance()

	x := ix.Get()
	y := iy.Get()
	for x != nil {
		// If y == nil, all of the remaining intervals in A are to the right of B,
		// so just yield them.
		if y == nil {
			intervals = append(intervals, x)

			ix.Advance()
			x = ix.Get()
			continue
		}
		// Split x into parts left and right of y.
		// The diagrams below show the bisection results for various situations.
		// if left.IsZero() && !right.IsZero()
		//             xxx
		// y1y1 y2y2 y3  y4y4
		//             xxx
		// or
		//   xxxxxxxxxxxx
		// y1y1 y2y2 y3  y4y4
		//
		// if !left.IsZero() && !right.IsZero()
		//       x1x1x1x1x1
		//         y1  y2
		//
		// if left.IsZero() && right.IsZero()
		//    x1x1x1x1  x2x2x2
		//  y1y1y1y1y1y1y1
		//
		// if !left.IsZero() && right.IsZero()
		//   x1x1  x2
		//     y1y1y1y1
		left, right := x.Bisect(y)

		// If the left side of x is non-zero, it can definitely be pushed to the
		// resulting interval set since no subsequent y value will intersect it.
		// The sequences look something like
		//         x1x1x1x1x1       OR   x1x1x1 x2
		//             y1 y2                       y1y1y1
		// left  = x1x1                  x1x1x1
		// right =       x1x1                            {zero}
		if !left.IsZero() {
			intervals = append(intervals, left)
		}

		if !right.IsZero() {
			// If the right side of x is non-zero:
			// 1) Right is the remaining portion of x that needs to be pushed.
			x = right
			// 2) It's not possible for current y to intersect it, so advance y. It's
			//    possible nextY() will intersect it, so don't push yet.
			iy.Advance()
			y = iy.Get()
		} else {
			// There's nothing left of x to push, so advance x.
			ix.Advance()
			x = ix.Get()
		}
	}

	// Setting s.intervals is the only side effect in this function.
	s.intervals = intervals
}

// Intersect destructively modifies the set by intersectin it with b.
func (s *Set) Intersect(b set.Set) {
	var ivals []interval.UnaryInterval
	set.Intersect(s, b, func(ival interval.UnaryInterval) bool {
		ivals = append(ivals, ival)
		return true
	})
	s.intervals = ivals
}

// Size returns number of intervals in this set
func (s *Set) Size() int {
	return len(s.intervals)
}

// String returns a human-friendly representation of the set.
func (s *Set) String() string {
	var strs []string
	for _, x := range s.intervals {
		strs = append(strs, fmt.Sprintf("%s", x))
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

// adjoinOrAppend adds an interval to the end of intervals unless that value
// directly adjoins the last element of intervals, in which case the last
// element will be replaced by the adjoined interval.
func adjoinOrAppend(intervals []interval.UnaryInterval, x interval.UnaryInterval) []interval.UnaryInterval {
	lastIndex := len(intervals) - 1
	if lastIndex == -1 {
		return append(intervals, x)
	}
	adjoined := intervals[lastIndex].Adjoin(x)
	if adjoined.IsZero() {
		return append(intervals, x)
	}
	intervals[lastIndex] = adjoined
	return intervals
}

// searchLow returns the first index in s.intervals that is not before x.
func (s *Set) searchLow(x interval.UnaryInterval) int {
	return sort.Search(len(s.intervals), func(i int) bool {
		return !s.intervals[i].Before(x)
	})
}

// searchLow returns the index of the first interval in s.intervals that is
// entirely after x.
func (s *Set) searchHigh(x interval.UnaryInterval) int {
	return sort.Search(len(s.intervals), func(i int) bool {
		return x.Before(s.intervals[i])
	})
}

// iterator returns a function that yields elements of the set in order.
func (s *Set) iterator(extents interval.UnaryInterval, forward bool) func() interval.UnaryInterval {
	low, high := s.searchLow(extents), s.searchHigh(extents)

	i, stride := low, 1
	if !forward {
		i, stride = high-1, -1
	}

	return func() interval.UnaryInterval {
		if i < 0 || i >= len(s.intervals) {
			return nil
		}
		x := s.intervals[i]
		i += stride
		return x
	}
}

// Between iterates over the intervals within extents set and calls f
// with each. If f returns false, iteration ceases.
//
// Any interval within the set that overlaps partially with extents is truncated
// before being passed to f.
func (s *Set) Between(extents interval.UnaryInterval) set.IntervalRange {
	// Begin = first index in s.intervals that is not before extents.
	begin := sort.Search(len(s.intervals), func(i int) bool {
		return !s.intervals[i].Before(extents)
	})

	return &intervalIterator{
		ivals: s.intervals[begin:],
		fn: func(interval interval.UnaryInterval) (interval.UnaryInterval, bool, bool) {
			// If the interval is after the extents, there will be no more overlap, so
			// break out of the loop.
			if extents.Before(interval) {
				return nil, false, false
			}
			portionOfInterval := extents.Intersect(interval)
			if portionOfInterval.IsZero() {
				return nil, true, false
			}
			return portionOfInterval, false, true
		},
	}
}

// All returns an ordered slice of all the intervals in the set.
func (s *Set) All() []interval.UnaryInterval {
	return append(make([]interval.UnaryInterval, 0, len(s.intervals)), s.intervals...)
}

// Freeze returns an immutable copy of this set.
func (s *Set) Freeze() *ImmutableSet {
	return NewImmutableSet(s.All()...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type intervalIterator struct {
	ivals []interval.UnaryInterval
	fn    func(interval.UnaryInterval) (interval.UnaryInterval, bool, bool)

	i   int
	cur interval.UnaryInterval
}

func (i *intervalIterator) GetSlice() []interval.UnaryInterval {
	var l []interval.UnaryInterval
	for i.Advance() {
		l = append(l, i.Get())
	}
	return l
}

func (i *intervalIterator) Iterator() set.IntervalIterator {
	return i
}

func (i *intervalIterator) Advance() (ok bool) {
	if i.i >= len(i.ivals) {
		i.cur = nil
		return false
	}
	retry := true
	for retry {
		i.cur, retry, ok = i.fn(i.ivals[i.i])
		if !ok {
			return ok
		}
	}
	i.i++
	return ok
}

func (i *intervalIterator) Get() interval.UnaryInterval {
	return i.cur
}
