package booking

import (
	"context"
	"time"

	"github.com/basgys/booking-consensys/app/iam"
	"github.com/deixis/errors"
	"github.com/deixis/pkg/utc"
	"github.com/deixis/spine/log"
)

type Service struct {
	rooms        *RoomsRepository
	reservations *ReservationRepository
}

func New(ctx context.Context) (*Service, error) {
	rooms, err := NewRoomsRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise room repository")
	}
	reservations, err := NewReservationRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise reservation repository")
	}

	return &Service{
		rooms:        rooms,
		reservations: reservations,
	}, nil
}

func (s *Service) ListRooms(
	ctx context.Context,
) ([]*Room, error) {
	return s.rooms.List(ctx)
}

func (s *Service) RoomsAvailabilities(
	ctx context.Context,
	from, to utc.UTC,
	d time.Duration,
) ([]*TimeInterval, error) {
	log.Trace(ctx, "booking.rooms.availabilities", "All Rooms availabilities",
		log.Stringer("from", from),
		log.Stringer("to", to),
		log.Stringer("duration", d),
	)
	return nil, errors.New("not implemented yet")
}

func (s *Service) RoomAvailabilities(
	ctx context.Context,
	roomRef string,
	from, to utc.UTC,
	d time.Duration,
) ([]*TimeInterval, error) {
	log.Trace(ctx, "booking.room.availabilities", "Room availabilities",
		log.String("roomRef", roomRef),
		log.Stringer("from", from),
		log.Stringer("to", to),
		log.Stringer("duration", d),
	)

	// TODO: Include duration
	return s.reservations.FreeRanges(ctx, roomRef, from, to)
}

func (s *Service) ListRoomReservations(
	ctx context.Context, roomRef string,
) ([]*Reservation, error) {
	return s.reservations.Reservations(ctx, roomRef)
}

func (s *Service) ReserveRoom(
	ctx context.Context,
	roomRef string,
	from utc.UTC,
	hours int64,
) (*Reservation, error) {
	acc, ok := iam.FromContext(ctx)
	if !ok {
		return nil, errors.PermissionDenied
	}

	res := Reservation{
		From:    from,
		To:      from.Add(time.Duration(hours) * time.Hour),
		RoomRef: roomRef,
		UserID:  acc.UserID,
	}
	if err := s.reservations.Reserve(ctx, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) CancelRoomReservation(
	ctx context.Context,
	roomRef string,
	id string,
) error {
	if _, ok := iam.FromContext(ctx); !ok {
		return errors.PermissionDenied
	}
	return s.reservations.Cancel(ctx, roomRef, id)
}
