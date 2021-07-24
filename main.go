package main

import (
	"context"
	"os"
	"strconv"

	"github.com/basgys/booking-consensys/app"
	"github.com/deixis/spine"
	"github.com/deixis/spine/net/http"
	"github.com/deixis/storage/kvdb"
	"github.com/deixis/storage/kvdb/driver/badger"
	"github.com/pkg/errors"
)

var (
	version = "dirty"
)

func main() {
	// Create spine
	block, err := spine.New("booking-api", nil)
	if err != nil {
		panic(errors.Wrap(err, "error initialising spine"))
	}
	block.Config().Version = version

	// Refer to spine instance as the main context
	var ctx context.Context
	ctx = block

	// Initialise storage
	os.Mkdir(".storage", 0770)
	store, err := badger.Open(".storage")
	if err != nil {
		panic(err)
	}
	ctx = kvdb.WithContext(ctx, store)

	// Initialises HTTP handler
	httpPort := parsePort(os.Getenv("HTTP_PORT"))
	httpServer := http.NewServer()
	block.RegisterService(&spine.ServiceRegistration{
		Name:   "http.booking-api",
		Host:   os.Getenv("IP"),
		Port:   httpPort,
		Server: httpServer,
		Tags:   []string{"http"},
	})

	// Init the application
	app, err := app.New(ctx)
	if err != nil {
		panic(errors.Wrap(err, "error initialising app"))
	}
	defer app.Close()

	// Initialise HTTP middlewares and endpoints
	app.HandleHTTP(httpServer)

	// Seed data
	// TODO: Only run in dev env
	if err := app.Seed(ctx); err != nil {
		panic(errors.Wrap(err, "error seeding"))
	}

	// Start serving requests
	if err := block.Serve(); err != nil {
		panic(err)
	}
}

func parsePort(s string) uint16 {
	p, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		panic(err)
	}
	return uint16(p)
}
