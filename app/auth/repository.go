package auth

import (
	"sync"
	"time"

	"github.com/basgys/booking-consensys/app/iam"
	"github.com/deixis/errors"
)

// ChallengePeriod defines how long a challenge remains valid
const ChallengePeriod = 3 * time.Minute

type ChallengeRepository struct {
	// Storing in memory for now to speed up development
	//
	// We currently store one challenge per address, which means we could potentially
	// prevent someone from logging in by spamming auth/challenge.
	// We could mitigate this issue by having multiple challenges per address,
	// but also with rate limiter.
	kv sync.Map
}

func (r *ChallengeRepository) Put(a iam.Address, challenge string) error {
	// TODO: Attach time to challenge
	r.kv.Store(a.String(), challenge)
	return nil
}

func (r *ChallengeRepository) Get(a iam.Address) (string, error) {
	// TODO: Return NotFound when TTL expired
	v, ok := r.kv.Load(a.String())
	if !ok {
		return "", errors.NotFound
	}
	return v.(string), nil
}
