package interval

// Interval is the interface for a continuous or discrete span. The interval is
// assumed to be inclusive of the starting point and exclusive of the ending
// point.
//
// All methods in the interface are non-destructive: Calls to the methods should
// not modify the interval. Furthermore, the implementation assumes an interval
// will not be mutated by user code, either.
type Interval interface {
	// Starting resturns the starting endpoint
	Starting() Endpoint
	// Ending resturns the ending endpoint
	Ending() Endpoint

	// After returns true when interval a `StartTouching` or is `EntirelyAfter`
	After(b Interval) bool
	// Before returns true when interval a is `EntirelyBefore` or `EndTouching`
	Before(b Interval) bool
	// Overlap returns whether the intervals overlap
	// Note: The relations `StartTouching` and `EndTouching` are not classified
	// as an overlap
	Overlap(b Interval) bool

	// IsZero returns true for the zero value of an interval.
	IsZero() bool
}

// UnaryInterval is a non-weighted `Interval`
type UnaryInterval interface {
	Interval

	// Intersect returns the intersection of an interval with another
	// interval. The function may panic if the other interval is incompatible.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===A===]
	//						[===B===]
	//					  |    |
	// Output:    [====]
	Intersect(UnaryInterval) UnaryInterval
	// Split returns two intervals, one on the lower side of x and one on the
	// upper side of x. The returned intervals are always within the range of the
	// original interval.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===============]
	//         |       X       |
	//         |       |       |
	// Output: [===L===]       |
	//				 |       [===R===]
	Split(x Endpoint) (left UnaryInterval, right UnaryInterval)
	// Bisect returns two intervals, one on the lower side of x and one on the
	// upper side of x, corresponding to the subtraction of x from the original
	// interval. The returned intervals are always within the range of the
	// original interval.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===============]
	//         |   [===X===]   |
	//         |   |       |   |
	// Output: [=L=]       [=R=]
	Bisect(x UnaryInterval) (left UnaryInterval, right UnaryInterval)
	// Adjoin returns the union of two intervals, if the intervals are exactly
	// adjacent, or the zero interval if they are not.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===A===]       |
	//         |       [===B===]
	//         |               |
	// Output: [===============]
	Adjoin(UnaryInterval) UnaryInterval
	// Encompass returns an interval that covers the exact extents of two
	// intervals.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [==A==]     [==B==]
	//         |                 |
	// Output: [=================]
	Encompass(UnaryInterval) UnaryInterval
}

// WeightedInterval is a weighted `Interval`
type WeightedInterval interface {
	Interval

	// Weight returns the interval weight
	Weight() uint64

	// Intersect returns the intersection of an interval with another
	// interval. The function may panic if the other interval is incompatible.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===A===]
	//						[===B===]
	//					  |    |
	// Output:    [====]
	Intersect(WeightedInterval) WeightedInterval
	// Add adds the given weight to the interval weight and returns
	// a new interval.
	Add(w uint64) WeightedInterval
	// Substract substracts the given weight from the interval weight and returns
	// a new interval.
	// When the a-b is a negative value or zero, IsZero will be true
	Substract(w uint64) WeightedInterval
	// Split returns two intervals, one on the lower side of x and one on the
	// upper side of x. The returned intervals are always within the range of the
	// original interval.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===============]
	//         |       X       |
	//         |       |       |
	// Output: [===L===]       |
	//				 |       [===R===]
	Split(x Endpoint) (left WeightedInterval, right WeightedInterval)
	// Bisect returns two intervals, one on the lower side of x and one on the
	// upper side of x, corresponding to the subtraction of x from the original
	// interval. The returned intervals are always within the range of the
	// original interval.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===============]
	//         |   [===X===]   |
	//         |   |       |   |
	// Output: [=L=]       [=R=]
	Bisect(x WeightedInterval) (left WeightedInterval, right WeightedInterval)
	// Adjoin returns the union of two intervals, if the intervals are exactly
	// adjacent, or the zero interval if they are not.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [===A===]       |
	//         |       [===B===]
	//         |               |
	// Output: [===============]
	Adjoin(WeightedInterval) WeightedInterval
	// Encompass returns an interval that covers the exact extents of two
	// intervals.
	//
	// e.g.
	// | | | | | | | | | | | | | |
	//  Input: [==A==]     [==B==]
	//         |                 |
	// Output: [=================]
	Encompass(WeightedInterval) WeightedInterval
}
