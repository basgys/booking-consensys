package timespan

import (
	"fmt"

	"github.com/basgys/booking-consensys/pkg/timeutil/interval"
	"github.com/basgys/booking-consensys/pkg/timeutil/set"
	"github.com/basgys/booking-consensys/pkg/timeutil/set/disjoint"
	"github.com/deixis/pkg/utc"
)

// Set is a finite set of time spans. Functions are provided for iterating over the
// spans and performing set operations (intersection, union, subtraction).
//
// This is a time span-specific implemention of interval.Set.
type Set struct {
	iset *disjoint.Set
}

// String returns a human readable version of the set.
func (s *Set) String() string {
	return s.iset.String()
}

// Empty returns a new, empty Set.
func Empty() *Set {
	return &Set{disjoint.Empty()}
}

// Copy returns a copy of a set that may be mutated without affecting the original.
func (s *Set) Copy() *Set {
	return &Set{s.iset.Copy()}
}

// Insert adds a single time span into the set.
func (s *Set) Insert(start, end utc.UTC) {
	if s.iset == nil {
		panic("timespan.Set not initialised")
	}
	if end < start {
		panic(fmt.Errorf("start %s before end %s", start, end))
	}
	s.iset.Insert(&Span{start, end})
}

// Remove removes a single time span from the set.
func (s *Set) Remove(start, end utc.UTC) {
	if s.iset == nil {
		panic("timespan.Set not initialised")
	}
	if end < start {
		panic(fmt.Errorf("start %s before end %s", start, end))
	}
	s.iset.Remove(&Span{start, end})
}

// Add performs an in-place union of two sets.
func (s *Set) Add(b *Set) {
	s.iset.Add(b.iset)
}

// Sub performs an in-place subtraction of set b from set a.
func (s *Set) Sub(b *Set) {
	s.iset.Sub(b.iset)
}

// Intersect performs an in-place intersection of sets a and b.
func (s *Set) Intersect(b *Set) {
	s.iset.Intersect(b.iset)
}

// Extent returns the start and end time that defines the entire timespan
// covering the set. The returned times are the zero value for an empty set.
func (s *Set) Extent() (utc.UTC, utc.UTC) {
	x := s.iset.Extent()
	if x == nil {
		return 0, 0
	}
	tr := tryOrPanic(x)
	return tr.Start, tr.End
}

// Empty reports if the extent of the set is zero.
func (s *Set) Empty() bool {
	x := s.iset.Extent()
	if x == nil {
		return true
	}
	return x.IsZero()
}

// Contains reports whether a time span is entirely contained within the set.
func (s *Set) Contains(start, end utc.UTC) bool {
	return s.iset.Contains(&Span{start, end})
}

// IntervalsBetween iterates over the time ranges within the set and calls f with the
// start (inclusive) and end (exclusive) of each. If f returns false, iteration
// ceases. Only intervals in between the provided start and end times are
// included.
func (s *Set) IntervalsBetween(extents interval.Interval) set.IntervalRange {
	tr := tryOrPanic(extents)
	return s.iset.Between(&Span{tr.Start, tr.End})
}

func (s *Set) Size() int {
	return s.iset.Size()
}

// IntervalReceiver is a function used for iterating over a set of time
// ranges. It takes the start and end times and returns true if the iteration
// should continue.
type IntervalReceiver func(start, end utc.UTC) bool

// ensureNonAdjoining returns a modified version of an IntervalsBetween callback
// function that will always be called with non-adjoining intervals. To do this,
// it returns a function that accumulates adjoining intervals, calling f with
// the combined interval. The second return value is a function that should be
// called after iteration is complete to ensure the last interval is sent to f.
func ensureNonAdjoining(f IntervalReceiver) (IntervalReceiver, func()) {
	last := &Span{}
	isDone := false
	doneFn := func() {
		if isDone {
			return
		}
		if !last.IsZero() {
			f(last.Start, last.End)
		}
	}
	receiveInterval := func(start, end utc.UTC) bool {
		if isDone {
			panic("should not be done")
		}
		current := &Span{start, end}
		adjoined := last.adjoin(current)
		if !adjoined.IsZero() {
			// Always continue if this interval adjoins the last one because the next
			// may also adjoin.
			last = adjoined
			return true
		}
		if !last.IsZero() {
			isDone = !f(last.Start, last.End)
			if isDone {
				return false //stop iteration
			}
		}
		last = current
		return true // continue iteration
	}
	return receiveInterval, doneFn
}
