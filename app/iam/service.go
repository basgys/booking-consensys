package iam

import (
	"context"

	"github.com/deixis/errors"
)

type Service struct {
	users  *UserRepository
	groups *GroupRepository
}

func New(ctx context.Context) (*Service, error) {
	users, err := NewUserRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error initialising user repository")
	}
	groups, err := NewGroupRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error initialising group repository")
	}

	return &Service{
		users:  users,
		groups: groups,
	}, nil
}
