package interval

// Rel is a relation in time between two intervals
type Rel int8

// Relation returns the interval relation of i in comparison to j.
// Note: It assumes that both intervals are valid
//
// e.g.
// | | | | | | | | | | | | | | | | | | | | |
//  	         |<=== j ===>|
//          	 |        	 |
//   <= i =>	 |           |						- EntirelyBefore
//   	  <= i =>|       	   |						- End Touching
//   	    <= i |=>		     |						- End Inside
//   	     <===|==== i ===>|			  		- Enclosing End Touching
//   	     <===|==== i ====|===>				- Enclosing
//   	         |<=== i ====|===>				- Enclosing Start Touching
//   	         |<=== i ===>|						- Exact Match
//   	         |<= i =>	   |						- Inside Start Touching
//   	         |  <= i =>	 |						- Inside
//   	         |    <= i =>|						- Inside End Touching
//   	         |         <=| i =>				- Start Inside
//   	         |           |<= i =>			- Start Touching
//   	         |           |	 <= i =>	- EntirelyAfter
func Relation(i, j Interval) Rel {
	switch i.Ending().Cmp(j.Starting()) {
	case -1:
		return EntirelyBefore
	case 0:
		return EndTouching
	}
	switch i.Starting().Cmp(j.Ending()) {
	case 0:
		return StartTouching
	case 1:
		return EntirelyAfter
	}

	x := i.Starting().Cmp(j.Starting())
	y := i.Ending().Cmp(j.Ending())
	switch x {
	case -1:
		switch y {
		case -1:
			return EndInside
		case 0:
			return EnclosingEndTouching
		case 1:
			return Enclosing
		}
	case 0:
		switch y {
		case -1:
			return InsideStartTouching
		case 0:
			return ExactMatch
		case 1:
			return EnclosingStartTouching
		}
	case 1:
		switch y {
		case -1:
			return Inside
		case 0:
			return InsideEndTouching
		case 1:
			return StartInside
		}
	}
	return UnknownRelation
}

// Overlap returns whether the relation represents two intervals that overlap
func (r Rel) Overlap() bool {
	return r >= EndInside && r <= StartInside
}

// Before returns true when interval a is entirely before or ends touching b.
func (r Rel) Before() bool {
	return r >= EntirelyBefore && r <= EndTouching
}

// After returns true when interval a starts touching or is entirely after b.
func (r Rel) After() bool {
	return r >= StartTouching && r <= EntirelyAfter
}

func (r Rel) String() string {
	switch r {
	case EntirelyBefore:
		return "EntirelyBefore"
	case EndTouching:
		return "EndTouching"
	case EndInside:
		return "EndInside"
	case EnclosingEndTouching:
		return "EnclosingEndTouching"
	case Enclosing:
		return "Enclosing"
	case EnclosingStartTouching:
		return "EnclosingStartTouching"
	case ExactMatch:
		return "ExactMatch"
	case InsideStartTouching:
		return "InsideStartTouching"
	case Inside:
		return "Inside"
	case InsideEndTouching:
		return "InsideEndTouching"
	case StartInside:
		return "StartInside"
	case StartTouching:
		return "StartTouching"
	case EntirelyAfter:
		return "EntirelyAfter"
	}
	return ""
}

const (
	UnknownRelation Rel = -1
	EntirelyBefore      = iota
	EndTouching
	EndInside
	EnclosingEndTouching
	Enclosing
	EnclosingStartTouching
	ExactMatch
	InsideStartTouching
	Inside
	InsideEndTouching
	StartInside
	StartTouching
	EntirelyAfter
)
