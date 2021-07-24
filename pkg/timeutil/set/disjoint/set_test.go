package disjoint_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/basgys/booking-consensys/pkg/timeutil/interval"
	"github.com/basgys/booking-consensys/pkg/timeutil/set"
	"github.com/basgys/booking-consensys/pkg/timeutil/set/disjoint"
)

func TestExtent(t *testing.T) {
	x := &span{20, 40}
	y := &span{60, 100}

	ival := disjoint.NewSet(x, y)
	if got, want := cast(ival.Extent()), (&span{20, 100}); !got.Equal(want) {
		t.Errorf("Extent() = %v, want %v", got, want)
	}
}

func TestAdd(t *testing.T) {
	x := disjoint.NewSet(&span{20, 40})
	y := disjoint.NewSet(&span{60, 111})

	if got, want := cast(x.Extent()), (&span{20, 40}); !got.Equal(want) {
		t.Errorf("Extent() = %v, want %v", got, want)
	}

	if got, want := cast(y.Extent()), (&span{60, 111}); !got.Equal(want) {
		t.Errorf("Extent() = %v, want %v", got, want)
	}

	x.Add(y)

	if got, want := cast(x.Extent()), (&span{20, 111}); !got.Equal(want) {
		t.Errorf("Extent() = %v, want %v", got, want)
	}

	for _, tt := range []struct {
		name string
		a    *disjoint.Set
		b    set.Set
		want []*span
	}{
		{
			"[20, 40) + [60,111)",
			disjoint.NewSet(&span{20, 40}),
			disjoint.NewSet(&span{60, 111}),
			[]*span{
				{20, 40},
				{60, 111},
			},
		},
		{
			"[20, 40) + [30,111) = [20, 111)",
			disjoint.NewSet(&span{20, 40}),
			disjoint.NewImmutableSet(&span{30, 111}),
			[]*span{
				{20, 111},
			},
		},
	} {
		u := disjoint.NewImmutableSet(tt.a.All()...).Union(tt.b)
		tt.a.Add(tt.b)
		if got := allIntervals(tt.a); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
		if got := allIntervals(u); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: [ImmutableSet] got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSub(t *testing.T) {
	for _, tt := range []struct {
		name string
		a    *disjoint.Set
		b    set.Set
		want []*span
	}{
		{
			"[20, 40) - [30,111)",
			disjoint.NewSet(&span{20, 40}),
			disjoint.NewSet(&span{30, 111}),
			[]*span{
				{20, 30},
			},
		},
		{
			"[0, 2) [4, 6) [8, 10) - [1, 2) [5, 6) [9, 10)   = [0, 1) [4, 5) [8, 9)",
			disjoint.NewSet(&span{0, 2}, &span{4, 6}, &span{8, 10}),
			disjoint.NewSet(&span{1, 2}, &span{5, 6}, &span{9, 10}),
			[]*span{{0, 1}, {4, 5}, {8, 9}},
		},
		{
			"[0...3)[10...13)...[90...93) - all odd numbers",
			func() *disjoint.Set {
				spans := []interval.UnaryInterval{}
				for i := 0; i < 100; i += 10 {
					spans = append(spans, &span{i, i + 3})
				}
				return disjoint.NewSet(spans...)
			}(),
			func() *disjoint.Set {
				spans := []interval.UnaryInterval{}
				for i := 1; i < 100; i += 2 {
					spans = append(spans, &span{i, i + 1})
				}
				return disjoint.NewSet(spans...)
			}(),
			func() []*span {
				spans := []*span{}
				for i := 0; i < 100; i += 10 {
					spans = append(spans, &span{i, i + 1}, &span{i + 2, i + 3})
				}
				return spans
			}(),
		},
	} {
		// Immutable set
		if got := allIntervals(tt.a.Freeze().Sub(tt.b)); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: [ImmutableSet] got %v, want %v", tt.name, got, tt.want)
		}

		// Mutable set
		tt.a.Sub(tt.b)
		if got := allIntervals(tt.a); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDisjointSet_Intersect(t *testing.T) {
	for _, tt := range []struct {
		name string
		a    *disjoint.Set
		b    set.Set
		want []*span
	}{
		{
			"[20, 40) intersect [30,111)",
			disjoint.NewSet(&span{20, 40}),
			disjoint.NewSet(&span{30, 111}),
			[]*span{{30, 40}},
		},
		{
			"[0, 2) [4, 6) [8, 10) intersect [1, 2) [5, 6) [9, 10)   = [1, 3) [5, 6) [9, 10)",
			disjoint.NewSet(&span{0, 2}, &span{4, 6}, &span{8, 10}),
			disjoint.NewSet(&span{1, 2}, &span{5, 6}, &span{9, 10}),
			[]*span{{1, 2}, {5, 6}, {9, 10}},
		},
		{
			"[0, 2) [5, 7) intersect [5, 7)   = [1, 2) [5, 6)",
			// [01...56...]
			// [.12345....]
			disjoint.NewSet(&span{0, 2}, &span{5, 7}),
			disjoint.NewSet(&span{1, 6}),
			[]*span{{1, 2}, {5, 6}},
		},
		{
			"[0, 2) [5, 7) intersect [5, 7)   = [1, 2) [5, 6)",
			disjoint.NewSet(&span{1, 6}),
			disjoint.NewSet(&span{0, 2}, &span{5, 7}),
			[]*span{{1, 2}, {5, 6}},
		},
		{
			"[0...7)[10...17)...[90...97) intersect (all odd numbers + {4, 14, ... 94})",
			func() *disjoint.Set {
				spans := []interval.UnaryInterval{}
				for i := 0; i < 100; i += 10 {
					spans = append(spans, &span{i, i + 7})
				}
				return disjoint.NewSet(spans...)
			}(),
			func() *disjoint.Set {
				spans := []interval.UnaryInterval{}
				for i := 0; i < 100; i += 10 {
					spans = append(
						spans,
						&span{i + 1, i + 2},
						&span{i + 3, i + 6},
						&span{i + 7, i + 8},
						&span{i + 9, i + 10},
					)
				}
				return disjoint.NewSet(spans...)
			}(),
			func() []*span {
				spans := []*span{}
				for i := 0; i < 100; i += 10 {
					spans = append(
						spans,
						&span{i + 1, i + 2},
						&span{i + 3, i + 6},
					)
				}
				return spans
			}(),
		},
	} {
		if got := allIntervals(tt.a.Freeze().Intersect(tt.b)); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: [ImmutableSet] got\n  %v, want\n  %v", tt.name, got, tt.want)
		}
		tt.a.Intersect(tt.b)
		if got := allIntervals(tt.a); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got\n  %v, want\n  %v", tt.name, got, tt.want)
		}
	}
}

func TestIntersect(t *testing.T) {
	for _, tt := range []struct {
		name string
		a    *disjoint.Set
		b    set.Set
		want []*span
	}{
		{
			"[20, 40) intersect [30,111)",
			disjoint.NewSet(&span{20, 40}),
			disjoint.NewSet(&span{30, 111}),
			[]*span{{30, 40}},
		},
		{
			"[0, 2) [4, 6) [8, 10) intersect [1, 2) [5, 6) [9, 10)   = [1, 3) [5, 6) [9, 10)",
			disjoint.NewSet(&span{0, 2}, &span{4, 6}, &span{8, 10}),
			disjoint.NewSet(&span{1, 2}, &span{5, 6}, &span{9, 10}),
			[]*span{{1, 2}, {5, 6}, {9, 10}},
		},
		{
			"[0, 2) [5, 7) intersect [5, 7)   = [1, 2) [5, 6)",
			// [01...56...]
			// [.12345....]
			disjoint.NewSet(&span{0, 2}, &span{5, 7}),
			disjoint.NewSet(&span{1, 6}),
			[]*span{{1, 2}, {5, 6}},
		},
		{
			"[0, 2) [5, 7) intersect [5, 7)   = [1, 2) [5, 6)",
			disjoint.NewSet(&span{1, 6}),
			disjoint.NewSet(&span{0, 2}, &span{5, 7}),
			[]*span{{1, 2}, {5, 6}},
		},
		{
			"[0...7)[10...17)...[90...97) intersect (all odd numbers + {4, 14, ... 94})",
			func() *disjoint.Set {
				spans := []interval.UnaryInterval{}
				for i := 0; i < 100; i += 10 {
					spans = append(spans, &span{i, i + 7})
				}
				return disjoint.NewSet(spans...)
			}(),
			func() *disjoint.Set {
				spans := []interval.UnaryInterval{}
				for i := 0; i < 100; i += 10 {
					spans = append(spans, &span{i + 1, i + 2}, &span{i + 3, i + 6}, &span{i + 7, i + 8}, &span{i + 9, i + 10})
				}
				return disjoint.NewSet(spans...)
			}(),
			func() []*span {
				spans := []*span{}
				for i := 0; i < 100; i += 10 {
					spans = append(spans, &span{i + 1, i + 2}, &span{i + 3, i + 6})
				}
				return spans
			}(),
		},
	} {
		var gotA []*span
		set.Intersect(tt.a, tt.b, func(ival interval.UnaryInterval) bool {
			gotA = append(gotA, ival.(*span))
			return true
		})
		if !reflect.DeepEqual(gotA, tt.want) {
			t.Errorf("%s: [Intersect] got\n  %v, want\n  %v", tt.name, gotA, tt.want)
		}

		var gotB []*span
		set.Intersect(tt.a, tt.b, func(ival interval.UnaryInterval) bool {
			gotB = append(gotB, ival.(*span))
			return true
		})
		if !reflect.DeepEqual(gotB, tt.want) {
			t.Errorf("%s: [Intersect] got\n  %v, want\n  %v", tt.name, gotB, tt.want)
		}
	}
}

func TestContains(t *testing.T) {
	for _, tt := range []struct {
		name string
		set  *disjoint.Set
		elem *span
		want bool
	}{
		{
			name: "{[0, 5), [10, 15)} contains [0, 5)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{0, 5},
			want: true,
		},
		{
			name: "{[0, 5), [10, 15)} does not contain [0, 6)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{0, 6},
			want: false,
		},
	} {
		if got := tt.set.Freeze().Contains(tt.elem); got != tt.want {
			t.Errorf("%s: [ImmutableSet] set.Contains(%s) = %t, want %t", tt.name, tt.elem, got, tt.want)
		}
		if got := tt.set.Contains(tt.elem); got != tt.want {
			t.Errorf("%s: set.Contains(%s) = %t, want %t", tt.name, tt.elem, got, tt.want)
		}
	}
}

func TestOverlap(t *testing.T) {
	for _, tt := range []struct {
		name string
		set  *disjoint.Set
		elem *span
		want bool
	}{
		{
			name: "{[0, 5), [10, 15)} overlaps [0, 5)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{0, 5},
			want: true,
		},
		{
			name: "{[0, 5), [10, 15)} overlaps [0, 6)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{0, 6},
			want: true,
		},
		{
			name: "{[0, 5), [10, 15)} overlaps [3, 11)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{3, 11},
			want: true,
		},
		{
			name: "{[0, 5), [10, 15)} overlaps [10, 15)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{10, 15},
			want: true,
		},
		{
			name: "{[0, 5), [10, 15)} overlaps [14, 16)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{14, 16},
			want: true,
		},
		{
			name: "{[0, 5), [10, 15)} does not overlaps [5, 10)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{5, 10},
			want: false,
		},
		{
			name: "{[0, 5), [10, 15)} does not overlaps [5, 10)]",
			set:  disjoint.NewSet(&span{0, 5}, &span{10, 15}),
			elem: &span{15, 20},
			want: false,
		},
		{
			name: "{[5, 10), [10, 15)} does not overlaps [0, 5)]",
			set:  disjoint.NewSet(&span{5, 10}, &span{10, 15}),
			elem: &span{0, 5},
			want: false,
		},
	} {
		if got := tt.set.Freeze().Overlap(tt.elem); got != tt.want {
			t.Errorf("%s: [ImmutableSet] set.Overlap(%s) = %t, want %t", tt.name, tt.elem, got, tt.want)
		}
		if got := tt.set.Overlap(tt.elem); got != tt.want {
			t.Errorf("%s: set.Overlap(%s) = %t, want %t", tt.name, tt.elem, got, tt.want)
		}
	}
}

func allIntervals(s set.Set) (spans []*span) {
	iter := s.Between(s.Extent()).Iterator()
	for iter.Advance() {
		spans = append(spans, cast(iter.Get()))
	}
	return spans
}

type span struct {
	Low, High int
}

func (s *span) String() string {
	return fmt.Sprintf("[%d, %d)", s.Low, s.High)
}

func (s *span) Equal(t *span) bool {
	return s.Low == t.Low && s.High == t.High
}

// Intersect returns the intersection of an disjoint with another
// disjoint. The function may panic if the other disjoint is incompatible.
func (s *span) Intersect(o interval.UnaryInterval) interval.UnaryInterval {
	t := cast(o)
	result := &span{
		max(s.Low, t.Low),
		min(s.High, t.High),
	}
	if result.Low < result.High {
		return result
	}
	return zero()
}

func (s *span) After(o interval.Interval) bool {
	return interval.Relation(s, o).After()
}

func (s *span) Before(o interval.Interval) bool {
	return interval.Relation(s, o).Before()
}

func (s *span) Overlap(o interval.Interval) bool {
	return interval.Relation(s, o).Overlap()
}

// IsZero returns true for the zero value of an disjoint.
func (s *span) IsZero() bool {
	return s.Low == 0 && s.High == 0
}

func (s *span) Split(x interval.Endpoint) (interval.UnaryInterval, interval.UnaryInterval) {
	if x.Cmp(s.Starting()) < 0 || x.Cmp(s.Ending()) > 0 {
		panic(fmt.Errorf(
			"S(%d,%d).Split(%d)",
			s.Starting(), s.Ending(), x,
		))
	}
	i, ok := x.(endpoint)
	if !ok {
		panic("incompatible endpoint")
	}
	maybeZero := func(Low, High int) *span {
		if Low == High {
			return zero()
		}
		return &span{Low, High}
	}
	return maybeZero(s.Low, int(i)), maybeZero(int(i), s.High)
}

// Bisect returns two disjoints, one on either lower side of x and one on the
// upper side of x, corresponding to the subtraction of x from the original
// disjoint. The returned disjoints are always within the range of the
// original disjoint.
func (s *span) Bisect(o interval.UnaryInterval) (interval.UnaryInterval, interval.UnaryInterval) {
	intersection := cast(s.Intersect(o))
	if intersection.IsZero() {
		if s.Before(o) {
			return s, zero()
		}
		return zero(), s
	}
	maybeZero := func(Low, High int) *span {
		if Low == High {
			return zero()
		}
		return &span{Low, High}
	}
	return maybeZero(s.Low, intersection.Low), maybeZero(intersection.High, s.High)

}

// Adjoin returns the union of two disjoints, if the disjoints are exactly
// adjacent, or the zero disjoint if they are not.
func (s *span) Adjoin(o interval.UnaryInterval) interval.UnaryInterval {
	t := cast(o)
	if s.High == t.Low {
		return &span{s.Low, t.High}
	}
	if t.High == s.Low {
		return &span{t.Low, s.High}
	}
	return zero()
}

// Encompass returns an disjoint that covers the exact extents of two
// disjoints.
func (s *span) Encompass(o interval.UnaryInterval) interval.UnaryInterval {
	t := cast(o)
	return &span{min(s.Low, t.Low), max(s.High, t.High)}
}

func (s *span) Starting() interval.Endpoint {
	return endpoint(s.Low)
}

func (s *span) Ending() interval.Endpoint {
	return endpoint(s.High)
}

type endpoint int

func (e endpoint) Cmp(b interval.Endpoint) int {
	other, ok := b.(endpoint)
	if !ok {
		panic("unsupported endpoint")
	}
	if e > other {
		return 1
	} else if e < other {
		return -1
	}
	return 0
}

// case returns a *span from an Interval interface, or it panics.
func cast(i interval.UnaryInterval) *span {
	x, ok := i.(*span)
	if !ok {
		panic(fmt.Errorf("disjoint must be an span: %v", i))
	}
	return x
}

// zero returns the zero value for span.
func zero() *span {
	return &span{}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
