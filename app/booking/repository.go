package booking

import (
	"bytes"
	"context"
	"encoding/gob"
	"strings"
	"time"

	"github.com/basgys/booking-consensys/pkg/timeutil/timespan"
	"github.com/deixis/errors"
	"github.com/deixis/pkg/utc"
	"github.com/deixis/storage/kvdb"
	"github.com/segmentio/ksuid"
)

var (
	firstKey kvdb.TupleElement
	lastKey  = kvdb.UUID{0xFF}
)

const (
	// minPrecision is the minimum precision applied to reservations.
	// That means everything is rounded up or down to the nearest precision
	minPrecision = 1 * time.Hour

	// minReservationDuration defines a minimum length for a room reservation.
	// All free ranges smaller than `minReservationDuration` will be discarded.
	minReservationDuration = 1 * time.Hour
)

type ReservationRepository struct {
	ss kvdb.Subspace
}

func NewReservationRepository(ctx context.Context) (*ReservationRepository, error) {
	store, ok := kvdb.FromContext(ctx)
	if !ok {
		return nil, kvdb.ErrNoConnectionFound
	}
	dir, err := store.CreateOrOpenDir([]string{"booking", "reservation"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open booking/reservation dir")
	}
	return &ReservationRepository{
		ss: dir,
	}, nil
}

func (r *ReservationRepository) Reservations(
	ctx context.Context, roomRef string,
) (reservations []*Reservation, err error) {
	roomRef = strings.ToUpper(roomRef)

	if roomRef == "" {
		return nil, errors.Bad(&errors.FieldViolation{
			Field:       "roomRef",
			Description: "Cannot query reservations without a room ref",
		})
	}

	_, err = kvdb.ReadTransact(ctx,
		func(ctx context.Context, tx kvdb.ReadTransaction) (interface{}, error) {
			data, err := tx.Get(r.ss.Pack([]kvdb.TupleElement{roomRef})).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound // room does not exist (yet)
			}
			err = gob.NewDecoder(bytes.NewReader(data)).Decode(&reservations)
			return nil, err
		},
	)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *ReservationRepository) Reserve(
	ctx context.Context, reservation *Reservation,
) error {
	reservation.From = reservation.From.Floor(minPrecision)
	reservation.To = reservation.To.Ceil(minPrecision)
	reservation.RoomRef = strings.TrimSpace(strings.ToUpper(reservation.RoomRef))

	if reservation.RoomRef == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "roomRef",
			Description: "A reservation requires a room reference to be valid",
		})
	}
	if reservation.From >= reservation.To {
		return errors.Bad(&errors.FieldViolation{
			Field:       "to",
			Description: "Invalid reservation interval",
		})
	}

	// Generate a K-Sortable Unique IDentifier based on the start date
	id, err := ksuid.NewRandomWithTime(reservation.From.Time())
	if err != nil {
		return err
	}
	reservation.ID = id.String()

	_, err = kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			roomKey := r.ss.Pack([]kvdb.TupleElement{reservation.RoomRef})
			var reservations []*Reservation
			// Load list
			data, err := tx.Get(roomKey).Get()
			if err != nil {
				return nil, err
			}
			if len(data) > 0 {
				rdr := bytes.NewReader(data)
				if err := gob.NewDecoder(rdr).Decode(&reservations); err != nil {
					return nil, errors.Wrap(err, "failed to unmarshal reservations")
				}
			}

			// Ensure it is in a free range
			timeset := timespan.Empty()
			for _, res := range reservations {
				timeset.Insert(res.From, res.To)
			}
			busy := timeset.IntervalsBetween(reservation.Interval()).Iterator().Advance()
			if busy {
				return nil, errors.Aborted(&errors.ConflictViolation{
					Resource:    "reservation",
					Description: "There is already a reservation on this range",
				})
			}

			// Add reservation to the list
			reservations = append(reservations, reservation)

			// Persist list
			var encoded bytes.Buffer
			if err := gob.NewEncoder(&encoded).Encode(reservations); err != nil {
				return nil, errors.Wrap(err, "failed to marshal reservations")
			}
			tx.Set(roomKey, encoded.Bytes())
			return nil, err
		},
	)
	return err
}

func (r *ReservationRepository) Cancel(
	ctx context.Context, roomRef string, reservationID string,
) error {
	roomRef = strings.TrimSpace(strings.ToUpper(roomRef))

	if roomRef == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "roomRef",
			Description: "A room reference is required to cancel a reservation",
		})
	}
	if reservationID == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "to",
			Description: "A reservation id is required to cancel a reservation",
		})
	}

	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			roomKey := r.ss.Pack([]kvdb.TupleElement{roomRef})
			// Load list
			var reservations []*Reservation
			data, err := tx.Get(roomKey).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound
			}
			rdr := bytes.NewReader(data)
			if err := gob.NewDecoder(rdr).Decode(&reservations); err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal reservations")
			}

			// FIXME: It is a highly inneficient way to to cancel a reservation
			var kept []*Reservation
			for _, res := range reservations {
				if res.ID != reservationID {
					kept = append(kept, res)
				}
			}
			if len(kept) == len(reservations) {
				return nil, errors.NotFound
			}

			// Persist list
			var encoded bytes.Buffer
			if err := gob.NewEncoder(&encoded).Encode(kept); err != nil {
				return nil, errors.Wrap(err, "failed to marshal reservations")
			}
			tx.Set(roomKey, encoded.Bytes())
			return nil, err
		},
	)
	return err
}

// FreeRanges returns a disjoint set of free ranges
func (r *ReservationRepository) FreeRanges(
	ctx context.Context, roomRef string, from, to utc.UTC,
) (ivals []*TimeInterval, err error) {
	roomRef = strings.ToUpper(roomRef)

	// Ensure interval is valid
	if to < from {
		return nil, errors.Bad(&errors.FieldViolation{
			Field:       "to",
			Description: "The interval is invalid. to is smaller than from",
		})
	}

	// Initialise an empty disjoint set
	timeset := timespan.Empty()

	// Set the whole range as available by default
	timeset.Insert(from, to)

	// Load reservations
	// Note: This step will close free ranges
	reservations, err := r.Reservations(ctx, roomRef)
	if err != nil {
		return nil, err
	}
	for _, r := range reservations {
		timeset.Remove(r.From, r.To)
	}

	// Convert to availabilities
	iter := timeset.IntervalsBetween(&timespan.Span{Start: from, End: to}).Iterator()
	for iter.Advance() {
		iv := iter.Get().(*timespan.Span)
		if iv.Duration() < minReservationDuration {
			continue
		}

		ivals = append(ivals, &TimeInterval{
			From: iv.Start,
			To:   iv.End,
		})
	}
	return ivals, nil
}

type RoomsRepository struct {
	ss kvdb.Subspace
}

func NewRoomsRepository(ctx context.Context) (*RoomsRepository, error) {
	store, ok := kvdb.FromContext(ctx)
	if !ok {
		return nil, kvdb.ErrNoConnectionFound
	}
	dir, err := store.CreateOrOpenDir([]string{"booking", "room"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open booking/room dir")
	}
	return &RoomsRepository{
		ss: dir,
	}, nil
}

func (r *RoomsRepository) List(ctx context.Context) (rooms []*Room, err error) {
	_, err = kvdb.ReadTransact(ctx,
		func(ctx context.Context, tx kvdb.ReadTransaction) (interface{}, error) {
			rng := kvdb.KeyRange{
				Begin: r.ss.Pack([]kvdb.TupleElement{firstKey}),
				End:   r.ss.Pack([]kvdb.TupleElement{lastKey}),
			}
			iter := tx.GetRange(rng).Iterator()
			for iter.Advance() {
				kv, err := iter.Get()
				if err != nil {
					return nil, err
				}
				room := Room{}
				err = gob.NewDecoder(bytes.NewReader(kv.Value)).Decode(&room)
				if err != nil {
					return nil, err
				}
				rooms = append(rooms, &room)
			}

			return nil, err
		},
	)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *RoomsRepository) Create(ctx context.Context, room *Room) error {
	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			key := r.ss.Pack([]kvdb.TupleElement{room.Ref})

			var encoded bytes.Buffer
			err := gob.NewEncoder(&encoded).Encode(&room)
			if err != nil {
				return nil, err
			}

			tx.Set(key, encoded.Bytes())
			return nil, err
		},
	)
	return err
}
