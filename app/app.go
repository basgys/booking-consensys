// Package app contains the app initialisation
// This has been moved to its own package in order to be reused in the
// integration tests
package app

import (
	"context"
	"fmt"
	"io"
	nethttp "net/http"
	"time"

	"github.com/basgys/booking-consensys/app/auth"
	"github.com/basgys/booking-consensys/app/booking"
	"github.com/basgys/booking-consensys/app/iam"
	"github.com/basgys/booking-consensys/pkg/jwtutil"
	"github.com/basgys/booking-consensys/pkg/mw"
	"github.com/deixis/errors"
	"github.com/deixis/pkg/utc"
	"github.com/deixis/spine"
	"github.com/deixis/spine/net/http"
	"github.com/deixis/storage/kvdb"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	ctx          context.Context
	store        kvdb.Store
	services     []interface{}
	httpHandlers []httpHandler
}

func New(ctx context.Context) (*App, error) {
	store, ok := kvdb.FromContext(ctx)
	if !ok {
		return nil, errors.New("cannot find KV store in context")
	}

	auths, err := auth.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error initialising auth service")
	}
	iams, err := iam.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error initialising iam service")
	}
	bookings, err := booking.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error initialising booking service")
	}

	return &App{
		ctx:   ctx,
		store: store,
		services: []interface{}{
			auths,
			iams,
			bookings,
		},
		httpHandlers: []httpHandler{
			auths,
			bookings,
		},
	}, nil
}

func (a *App) Seed(ctx context.Context) error {
	accounts, err := iam.NewAccountRepository(ctx)
	if err != nil {
		return errors.Wrap(err, "error initialising account repository")
	}
	users, err := iam.NewUserRepository(ctx)
	if err != nil {
		return errors.Wrap(err, "error initialising user repository")
	}
	groups, err := iam.NewGroupRepository(ctx)
	if err != nil {
		return errors.Wrap(err, "error initialising group repository")
	}
	rooms, err := booking.NewRoomsRepository(ctx)
	if err != nil {
		return errors.Wrap(err, "error initialising room repository")
	}
	reservations, err := booking.NewReservationRepository(ctx)
	if err != nil {
		return errors.Wrap(err, "error initialising reservation repository")
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate private key for seed user")
	}
	addr := iam.Address(crypto.PubkeyToAddress(privKey.PublicKey))

	// Create groups
	coke := &iam.Group{
		ID:  uuid.New().String(),
		Ref: "COKE",
	}
	groups.Create(ctx, coke)
	pepsi := &iam.Group{
		ID:  uuid.New().String(),
		Ref: "PEPSI",
	}
	groups.Create(ctx, pepsi)

	// Create user with its account
	usr := &iam.User{
		ID:      uuid.New().String(),
		GroupID: coke.ID,
	}
	users.Create(ctx, usr)
	acc := &iam.Account{
		Address: iam.Address(addr),
		UserID:  usr.ID,
	}
	accounts.Create(ctx, acc)

	// Create some random reservations
	for _, g := range []string{"C", "P"} {
		for i := 1; i <= 10; i++ {
			roomRef := fmt.Sprintf("%s%02d", g, i)
			rooms.Create(ctx, &booking.Room{
				Ref: roomRef,
			})

			res := booking.Reservation{
				RoomRef: roomRef,
				From:    utc.MustParse("2021-07-26T12:00:00Z"),
				To:      utc.MustParse("2021-07-26T13:00:00Z"),
				UserID:  usr.ID,
			}
			err := reservations.Reserve(ctx, &res)
			switch {
			case err == nil || errors.IsAborted(err):
			default:
				return errors.Wrapf(err, "failed to create reservation for room %s", roomRef)
			}
		}
	}

	// Create session
	token := jwt.NewWithClaims(jwtutil.StandardMethod, jwt.StandardClaims{
		Id:        acc.ID(),
		ExpiresAt: int64(utc.Now().Add(30 * 24 * time.Hour)),
	})
	signedToken, err := token.SignedString([]byte("a-hardcoded-secret-is-the-safest-secret"))
	if err != nil {
		return errors.Wrap(err, "failed to sign JWT")
	}

	fmt.Println("====================================")
	fmt.Println("Seed data")
	fmt.Println("")
	fmt.Println("JWT:", signedToken)
	fmt.Println("Account:", addr.String())
	fmt.Println("User:", usr.ID)
	fmt.Println("")
	fmt.Println("Groups:")
	fmt.Println("Coke:", coke.ID)
	fmt.Println("Pepsi:", pepsi.ID)
	fmt.Println("====================================")
	return nil
}

func (a *App) Close() error {
	for _, svc := range a.services {
		if c, ok := svc.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *App) HandleHTTP(srv *http.Server) {
	ni := &mw.NodeInfo{
		Name:    "",
		Version: "",
	}
	if app, ok := a.ctx.(*spine.App); ok {
		ni.Name = app.Config().Node
		ni.Version = app.Config().Version
	}
	store := mw.Store{
		S: a.store,
	}

	// Add middlewares
	// Create HTTP handler with middlewares
	srv.Append(ni.ReturnNodeInfo)
	srv.Append(store.Inject)

	// Return 200 OK on / for load balancer health check
	srv.HandleFunc("/", http.GET, httpOK)

	for _, h := range a.httpHandlers {
		h.HandleHTTP(srv)
	}
	srv.HandleEndpoint(&matchAll{})
}

type httpHandler interface {
	HandleHTTP(srv *http.Server)
}

func httpOK(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	w.Head(http.StatusOK)
}

type matchAll struct{}

func (m *matchAll) Path() string   { return "/" }
func (m *matchAll) Method() string { return http.GET }

func (m *matchAll) Attach(r *mux.Router, f func(nethttp.ResponseWriter, *nethttp.Request)) {
	r.PathPrefix("/").HandlerFunc(f)
}

func (m *matchAll) Serve(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Head(http.StatusNotFound)
}
