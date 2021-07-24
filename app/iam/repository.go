package iam

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/deixis/errors"
	"github.com/deixis/storage/kvdb"
)

type GroupRepository struct {
	ss kvdb.Subspace
}

func NewGroupRepository(ctx context.Context) (*GroupRepository, error) {
	store, ok := kvdb.FromContext(ctx)
	if !ok {
		return nil, kvdb.ErrNoConnectionFound
	}
	dir, err := store.CreateOrOpenDir([]string{"iam", "group"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open iam/group dir")
	}
	return &GroupRepository{
		ss: dir,
	}, nil
}

func (r *GroupRepository) Create(
	ctx context.Context, g *Group,
) error {
	if g.ID == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing group ID",
		})
	}

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(g); err != nil {
		return errors.Wrap(err, "failed to marshal group")
	}

	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			key := r.ss.Pack([]kvdb.TupleElement{g.ID})
			data, err := tx.Get(key).Get()
			if err != nil {
				return nil, err
			}
			if len(data) > 0 {
				return nil, errors.Aborted(&errors.ConflictViolation{
					Resource:    "group:" + g.ID,
					Description: "Group has already been created",
				})
			}

			tx.Set(key, encoded.Bytes())
			return nil, nil
		},
	)
	return err
}

func (r *GroupRepository) Get(
	ctx context.Context, id string,
) (*Group, error) {
	v, err := kvdb.ReadTransact(ctx,
		func(ctx context.Context, tx kvdb.ReadTransaction) (interface{}, error) {
			data, err := tx.Get(r.ss.Pack([]kvdb.TupleElement{id})).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound
			}
			return data, nil
		},
	)
	if err != nil {
		return nil, err
	}

	g := &Group{}
	if err := gob.NewDecoder(bytes.NewReader(v.([]byte))).Decode(g); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal group")
	}
	return g, nil
}

func (r *GroupRepository) Update(
	ctx context.Context, g *Group,
) error {
	if g.ID == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing group ID",
		})
	}

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(g); err != nil {
		return errors.Wrap(err, "failed to marshal group")
	}

	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			key := r.ss.Pack([]kvdb.TupleElement{g.ID})
			data, err := tx.Get(key).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound
			}

			tx.Set(key, encoded.Bytes())
			return nil, nil
		},
	)
	return err
}

func (r *GroupRepository) Delete(
	ctx context.Context, id string,
) error {
	if id == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing group ID",
		})
	}
	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			tx.Clear(r.ss.Pack([]kvdb.TupleElement{id}))
			return nil, nil
		},
	)
	return err
}

type UserRepository struct {
	ss kvdb.Subspace
}

func NewUserRepository(ctx context.Context) (*UserRepository, error) {
	store, ok := kvdb.FromContext(ctx)
	if !ok {
		return nil, kvdb.ErrNoConnectionFound
	}
	dir, err := store.CreateOrOpenDir([]string{"iam", "user"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open iam/user dir")
	}
	return &UserRepository{
		ss: dir,
	}, nil
}

func (r *UserRepository) Create(
	ctx context.Context, u *User,
) error {
	if u.ID == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing user ID",
		})
	}

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(u); err != nil {
		return errors.Wrap(err, "failed to marshal user")
	}

	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			key := r.ss.Pack([]kvdb.TupleElement{u.ID})
			data, err := tx.Get(key).Get()
			if err != nil {
				return nil, err
			}
			if len(data) > 0 {
				return nil, errors.Aborted(&errors.ConflictViolation{
					Resource:    "user:" + u.ID,
					Description: "User has already been created",
				})
			}

			tx.Set(key, encoded.Bytes())
			return nil, nil
		},
	)
	return err
}

func (r *UserRepository) Get(
	ctx context.Context, id string,
) (*User, error) {
	v, err := kvdb.ReadTransact(ctx,
		func(ctx context.Context, tx kvdb.ReadTransaction) (interface{}, error) {
			data, err := tx.Get(r.ss.Pack([]kvdb.TupleElement{id})).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound
			}
			return data, nil
		},
	)
	if err != nil {
		return nil, err
	}

	u := &User{}
	if err := gob.NewDecoder(bytes.NewReader(v.([]byte))).Decode(u); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user")
	}
	return u, nil
}

func (r *UserRepository) Update(
	ctx context.Context, u *User,
) error {
	if u.ID == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing user ID",
		})
	}

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(u); err != nil {
		return errors.Wrap(err, "failed to marshal user")
	}

	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			key := r.ss.Pack([]kvdb.TupleElement{u.ID})
			data, err := tx.Get(key).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound
			}

			tx.Set(key, encoded.Bytes())
			return nil, nil
		},
	)
	return err
}

func (r *UserRepository) Delete(
	ctx context.Context, id string,
) error {
	if id == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing user ID",
		})
	}
	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			tx.Clear(r.ss.Pack([]kvdb.TupleElement{id}))
			return nil, nil
		},
	)
	return err
}

type AccountRepository struct {
	ss kvdb.Subspace
}

func NewAccountRepository(ctx context.Context) (*AccountRepository, error) {
	store, ok := kvdb.FromContext(ctx)
	if !ok {
		return nil, kvdb.ErrNoConnectionFound
	}
	dir, err := store.CreateOrOpenDir([]string{"iam", "account"})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open iam/account dir")
	}
	return &AccountRepository{
		ss: dir,
	}, nil
}

func (r *AccountRepository) Create(
	ctx context.Context, a *Account,
) error {
	id := a.Address.String()
	if id == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing account ID",
		})
	}

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(a); err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}

	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			key := r.ss.Pack([]kvdb.TupleElement{id})
			data, err := tx.Get(key).Get()
			if err != nil {
				return nil, err
			}
			if len(data) > 0 {
				return nil, errors.Aborted(&errors.ConflictViolation{
					Resource:    "account:" + id,
					Description: "Account has already been created",
				})
			}

			tx.Set(key, encoded.Bytes())
			return nil, nil
		},
	)
	return err
}

func (r *AccountRepository) Get(
	ctx context.Context, addr Address,
) (*Account, error) {
	id := addr.String()
	v, err := kvdb.ReadTransact(ctx,
		func(ctx context.Context, tx kvdb.ReadTransaction) (interface{}, error) {
			data, err := tx.Get(r.ss.Pack([]kvdb.TupleElement{id})).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound
			}
			return data, nil
		},
	)
	if err != nil {
		return nil, err
	}

	u := &Account{}
	if err := gob.NewDecoder(bytes.NewReader(v.([]byte))).Decode(u); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account")
	}
	return u, nil
}

func (r *AccountRepository) Update(
	ctx context.Context, a *Account,
) error {
	id := a.Address.String()
	if id == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing account ID",
		})
	}

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(a); err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}

	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			key := r.ss.Pack([]kvdb.TupleElement{id})
			data, err := tx.Get(key).Get()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, errors.NotFound
			}

			tx.Set(key, encoded.Bytes())
			return nil, nil
		},
	)
	return err
}

func (r *AccountRepository) Delete(
	ctx context.Context, addr Address,
) error {
	id := addr.String()
	if id == "" {
		return errors.Bad(&errors.FieldViolation{
			Field:       "id",
			Description: "Missing account ID",
		})
	}
	_, err := kvdb.Transact(ctx,
		func(ctx context.Context, tx kvdb.Transaction) (interface{}, error) {
			tx.Clear(r.ss.Pack([]kvdb.TupleElement{id}))
			return nil, nil
		},
	)
	return err
}
