package interval

// An Endpoint is a point in time marking the beginning/end of an Interval.
type Endpoint interface {
	// Cmp compares the endpoints and returns:
	//
	//     -1 if a <  b
	//      0 if a == b
	//     +1 if a >  b
	//
	Cmp(b Endpoint) int
}
