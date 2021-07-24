package booking_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/basgys/booking-consensys/app/booking"
	"github.com/deixis/errors"
	"github.com/deixis/pkg/utc"
	"github.com/deixis/storage/kvdb"
	"github.com/deixis/storage/kvdb/driver/badger"
)

const (
	storageFolder = ".storage"
)

func TestMain(m *testing.M) {
	var code int
	func() {
		// Prepare storage folder
		os.Mkdir(storageFolder, 0770)
		defer os.RemoveAll(storageFolder)

		code = m.Run()
	}()

	os.Exit(code)
}

// TestReservation_Reserve test a reservation creation on a new room
func TestReservation_Reserve(t *testing.T) {
	ctx, err := loadStorage(t.Name())
	if err != nil {
		t.Fatal("error loading storage", err)
	}

	reservations, err := booking.NewReservationRepository(ctx)
	if err != nil {
		t.Fatal("error opening repository", err)
	}

	res := &booking.Reservation{
		RoomRef: "C01",
		From:    utc.MustParse("2021-08-01T12:00:00Z"),
		To:      utc.MustParse("2021-08-01T13:00:00Z"),
		UserID:  "foo",
	}
	if err := reservations.Reserve(ctx, res); err != nil {
		t.Error("expect to create a reservation, but got", err)
	}
}

func TestReservation_Conflict(t *testing.T) {
	ctx, err := loadStorage(t.Name())
	if err != nil {
		t.Fatal("error loading storage", err)
	}

	reservations, err := booking.NewReservationRepository(ctx)
	if err != nil {
		t.Fatal("error opening repository", err)
	}

	res := &booking.Reservation{
		RoomRef: "C01",
		From:    utc.MustParse("2021-08-01T12:00:00Z"),
		To:      utc.MustParse("2021-08-01T13:00:00Z"),
		UserID:  "foo",
	}
	if err := reservations.Reserve(ctx, res); err != nil {
		t.Error("expect to create a reservation, but got", err)
	}
	res = &booking.Reservation{
		RoomRef: "C01",
		From:    utc.MustParse("2021-08-01T12:00:00Z"),
		To:      utc.MustParse("2021-08-01T13:00:00Z"),
		UserID:  "bar",
	}
	err = reservations.Reserve(ctx, res)
	if !errors.IsAborted(err) {
		t.Error("expect to get a conflict, but got", err)
	}
}

func TestReservation_Cancellation(t *testing.T) {
	ctx, err := loadStorage(t.Name())
	if err != nil {
		t.Fatal("error loading storage", err)
	}

	reservations, err := booking.NewReservationRepository(ctx)
	if err != nil {
		t.Fatal("error opening repository", err)
	}

	res := &booking.Reservation{
		RoomRef: "C01",
		From:    utc.MustParse("2021-08-01T12:00:00Z"),
		To:      utc.MustParse("2021-08-01T13:00:00Z"),
		UserID:  "foo",
	}
	if err := reservations.Reserve(ctx, res); err != nil {
		t.Error("expect to create a reservation, but got", err)
	}
	if err := reservations.Cancel(ctx, res.RoomRef, res.ID); err != nil {
		t.Error("expect to cancel a reservation, but got", err)
	}
	if err := reservations.Reserve(ctx, res); err != nil {
		t.Error("expect a room to be free after a cancellation, but got", err)
	}
}

// TestReservation_RoomIsolation ensures that rooms don't share the same
// schedule
func TestReservation_RoomIsolation(t *testing.T) {
	ctx, err := loadStorage(t.Name())
	if err != nil {
		t.Fatal("error loading storage", err)
	}

	reservations, err := booking.NewReservationRepository(ctx)
	if err != nil {
		t.Fatal("error opening repository", err)
	}

	res := &booking.Reservation{
		RoomRef: "C01",
		From:    utc.MustParse("2021-08-01T12:00:00Z"),
		To:      utc.MustParse("2021-08-01T13:00:00Z"),
		UserID:  "foo",
	}
	if err := reservations.Reserve(ctx, res); err != nil {
		t.Error("expect to create a reservation, but got", err)
	}
	// Same time, different room
	res = &booking.Reservation{
		RoomRef: "C02",
		From:    utc.MustParse("2021-08-01T12:00:00Z"),
		To:      utc.MustParse("2021-08-01T13:00:00Z"),
		UserID:  "foo",
	}
	if err := reservations.Reserve(ctx, res); err != nil {
		t.Error("expect to rooms to be isolated, but got", err)
	}
}

func loadStorage(name string) (context.Context, error) {
	store, err := badger.Open(path.Join(storageFolder, name))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ctx = kvdb.WithContext(ctx, store)
	return ctx, nil
}
