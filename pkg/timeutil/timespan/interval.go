package timespan

import (
	"fmt"
	"time"

	"github.com/basgys/booking-consensys/pkg/timeutil/interval"
	"github.com/deixis/pkg/utc"
)

// Span is an implementation of interval.Interval.
type Span struct {
	Start utc.UTC
	End   utc.UTC
}

// Duration returns the interval duration
func (ts *Span) Duration() time.Duration {
	return ts.End.Distance(ts.Start)
}

func (ts *Span) String() string {
	return fmt.Sprintf("[%s, %s)", ts.Start, ts.End)
}

func (ts *Span) Equal(b *Span) bool {
	return ts.Start == b.Start && ts.End == b.End
}

func (ts *Span) intersect(b *Span) *Span {
	result := &Span{
		max(ts.Start, b.Start),
		min(ts.End, b.End),
	}
	if result.Start < result.End {
		return result
	}
	return &Span{}
}

func (ts *Span) IsZero() bool {
	return ts.Start.IsZero() && ts.End.IsZero()
}

func (ts *Span) After(o interval.Interval) bool {
	return interval.Relation(ts, o).After()
}

func (ts *Span) Before(o interval.Interval) bool {
	return interval.Relation(ts, o).Before()
}

func (ts *Span) Overlap(o interval.Interval) bool {
	return interval.Relation(ts, o).Overlap()
}

func (ts *Span) Split(x interval.Endpoint) (interval.UnaryInterval, interval.UnaryInterval) {
	// t := tryOrPanic(other)
	panic("implement me")
}

func (ts *Span) Bisect(other interval.UnaryInterval) (interval.UnaryInterval, interval.UnaryInterval) {
	b := tryOrPanic(other)
	intersection := ts.intersect(b)
	if intersection.IsZero() {
		if ts.Before(b) {
			return ts, &Span{}
		}
		return &Span{}, ts
	}
	maybeZero := func(s, e utc.UTC) *Span {
		if s == e {
			return &Span{}
		}
		return &Span{s, e}
	}
	return maybeZero(ts.Start, intersection.Start), maybeZero(intersection.End, ts.End)
}

func (ts *Span) Intersect(other interval.UnaryInterval) interval.UnaryInterval {
	return ts.intersect(tryOrPanic(other))
}

func (ts *Span) Adjoin(other interval.UnaryInterval) interval.UnaryInterval {
	return ts.adjoin(tryOrPanic(other))
}

func (ts *Span) Encompass(other interval.UnaryInterval) interval.UnaryInterval {
	return ts.encompass(tryOrPanic(other))
}

func (ts *Span) Starting() interval.Endpoint {
	return endpoint(ts.Start)
}

func (ts *Span) Ending() interval.Endpoint {
	return endpoint(ts.End)
}

func (ts *Span) adjoin(b *Span) *Span {
	if ts.End == b.Start {
		return &Span{ts.Start, b.End}
	}
	if b.End == ts.Start {
		return &Span{b.Start, ts.End}
	}
	return &Span{}
}

func (ts *Span) encompass(b *Span) interval.UnaryInterval {
	return &Span{min(ts.Start, b.Start), max(ts.End, b.End)}
}

type endpoint utc.UTC

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

func max(a, b utc.UTC) utc.UTC {
	if a > b {
		return a
	}
	return b
}

func min(a, b utc.UTC) utc.UTC {
	if a < b {
		return a
	}
	return b
}

func tryOrPanic(i interval.Interval) *Span {
	tr, ok := i.(*Span)
	if !ok {
		panic(fmt.Errorf("interval must be a time range: %v", i))
	}
	return tr
}
