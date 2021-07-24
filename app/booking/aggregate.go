package booking

import (
	"time"

	"github.com/basgys/booking-consensys/pkg/timeutil/interval"
	"github.com/basgys/booking-consensys/pkg/timeutil/timespan"
	"github.com/deixis/pkg/utc"
)

type Room struct {
	Ref string `json:"ref"`
}

type Reservation struct {
	ID      string  `json:"id"`
	From    utc.UTC `json:"from"`
	To      utc.UTC `json:"to"`
	RoomRef string  `json:"roomRef"`
	UserID  string  `json:"userId"`
}

func (r *Reservation) Interval() interval.Interval {
	return &timespan.Span{Start: r.From, End: r.To}
}

// TimeInterval represents a contiguous range of time periods
type TimeInterval struct {
	From utc.UTC `json:"from"`
	To   utc.UTC `json:"to"`
}

// Duration returns the distance in time between From and To
func (i *TimeInterval) Duration() time.Duration {
	return i.To.Distance(i.From)
}
