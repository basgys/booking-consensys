package iam_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/basgys/booking-consensys/app/iam"
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

// TestGroup_CRUD calls each CRUD operation to make sure nothing returns
// an error.
//
// TODO: This test does not cover much and should be rewritten
func TestGroup_CRUD(t *testing.T) {
	ctx, err := loadStorage(t.Name())
	if err != nil {
		t.Fatal("error loading storage", err)
	}

	groups, err := iam.NewGroupRepository(ctx)
	if err != nil {
		t.Fatal("error opening repository", err)
	}

	g := &iam.Group{
		ID:  "coke-123",
		Ref: "Coke",
	}
	if err := groups.Create(ctx, g); err != nil {
		t.Fatal("error creating group", err)
	}
	if err := groups.Update(ctx, g); err != nil {
		t.Fatal("error updating group", err)
	}
	g, err = groups.Get(ctx, g.ID)
	if err != nil {
		t.Fatal("error loading group", err)
	}
	if err := groups.Delete(ctx, g.ID); err != nil {
		t.Fatal("error deleting group", err)
	}
}

// TestGroup_CRUD calls each CRUD operation to make sure nothing returns
// an error.
//
// TODO: This test does not cover much and should be rewritten
func TestUser_CRUD(t *testing.T) {
	ctx, err := loadStorage(t.Name())
	if err != nil {
		t.Fatal("error loading storage", err)
	}

	users, err := iam.NewUserRepository(ctx)
	if err != nil {
		t.Fatal("error opening repository", err)
	}

	u := &iam.User{
		ID: "coke-123",
	}
	if err := users.Create(ctx, u); err != nil {
		t.Fatal("error creating user", err)
	}
	if err := users.Update(ctx, u); err != nil {
		t.Fatal("error updating user", err)
	}
	u, err = users.Get(ctx, u.ID)
	if err != nil {
		t.Fatal("error loading user", err)
	}
	if err := users.Delete(ctx, u.ID); err != nil {
		t.Fatal("error deleting user", err)
	}
}

// TestGroup_CRUD calls each CRUD operation to make sure nothing returns
// an error.
//
// TODO: This test does not cover much and should be rewritten
func TestAccount_CRUD(t *testing.T) {
	ctx, err := loadStorage(t.Name())
	if err != nil {
		t.Fatal("error loading storage", err)
	}

	accounts, err := iam.NewAccountRepository(ctx)
	if err != nil {
		t.Fatal("error opening repository", err)
	}

	addr, err := iam.ParseAddress("0x35F659Ec81bb9A38ae576140107a6c5C8AE55900")
	if err != nil {
		t.Fatal("failed to parse address", err)
	}
	acc := &iam.Account{
		Address: addr,
	}
	if err := accounts.Create(ctx, acc); err != nil {
		t.Fatal("error creating account", err)
	}
	if err := accounts.Update(ctx, acc); err != nil {
		t.Fatal("error updating account", err)
	}
	acc, err = accounts.Get(ctx, acc.Address)
	if err != nil {
		t.Fatal("error loading account", err)
	}
	if err := accounts.Delete(ctx, acc.Address); err != nil {
		t.Fatal("error deleting account", err)
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
