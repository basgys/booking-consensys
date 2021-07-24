package disjoint

import (
	"github.com/basgys/booking-consensys/pkg/timeutil/interval"
	"github.com/basgys/booking-consensys/pkg/timeutil/set"
)

// ImmutableSet is a set of interval objects. It provides various set theory
// operations.
type ImmutableSet struct {
	set *Set
}

// NewImmutableSet returns a new set given a sorted slice of intervals. This
// function panics if the intervals are not sorted.
func NewImmutableSet(intervals ...interval.UnaryInterval) *ImmutableSet {
	return &ImmutableSet{NewSet(intervals...)}
}

// Extent returns the Interval defined by the minimum and maximum values of the
// set.
func (s *ImmutableSet) Extent() interval.UnaryInterval {
	return s.set.Extent()
}

// Contains reports whether an interval is entirely contained by the set.
func (s *ImmutableSet) Contains(ival interval.UnaryInterval) bool {
	return s.set.Contains(ival)
}

// Union returns a set with the contents of this set and another set.
func (s *ImmutableSet) Union(b set.Set) *ImmutableSet {
	union := s.set.Copy()
	union.Add(b)
	return &ImmutableSet{union}
}

// Sub returns a set without the intervals of another set.
func (s *ImmutableSet) Sub(b set.Set) *ImmutableSet {
	x := s.set.Copy()
	x.Sub(b)
	return &ImmutableSet{x}
}

// Intersect returns the intersection of two sets.
func (s *ImmutableSet) Intersect(b set.Set) *ImmutableSet {
	x := s.set.Copy()
	x.Intersect(b)
	return &ImmutableSet{x}
}

func (s *ImmutableSet) Overlap(ival interval.UnaryInterval) bool {
	return s.set.Overlap(ival)
}

func (s *ImmutableSet) Between(extents interval.UnaryInterval) set.IntervalRange {
	return s.set.Between(extents)
}

func (s *ImmutableSet) All() []interval.UnaryInterval {
	return s.set.All()
}

func (s *ImmutableSet) Size() int {
	return s.set.Size()
}

// String returns a human-friendly representation of the set.
func (s *ImmutableSet) String() string {
	return s.set.String()
}
