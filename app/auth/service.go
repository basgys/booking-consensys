package auth

import (
	"context"
	"time"

	"github.com/basgys/booking-consensys/app/auth/ethereum"
	"github.com/basgys/booking-consensys/app/iam"
	"github.com/basgys/booking-consensys/pkg/jwtutil"
	"github.com/deixis/errors"
	"github.com/deixis/pkg/utc"
	"github.com/deixis/spine/log"
	"github.com/form3tech-oss/jwt-go"
)

const (
	sessionTTL = 2 * time.Hour
)

type Service struct {
	Secret string

	auths      Auth
	challenges ChallengeRepository
	accounts   *iam.AccountRepository
}

func New(ctx context.Context) (*Service, error) {
	accounts, err := iam.NewAccountRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise account repository")
	}

	return &Service{
		// FIXME: Add secret to configuration file
		Secret:   "a-hardcoded-secret-is-the-safest-secret",
		auths:    &ethereum.Auth{},
		accounts: accounts,
	}, nil
}

func (s *Service) Challenge(ctx context.Context, address *iam.Address) (string, error) {
	if address == nil {
		return "", errors.Bad(&errors.FieldViolation{
			Field:       "account",
			Description: "Missing account id",
		})
	}

	log.Trace(ctx, "auth.challenge", "Challenge",
		log.Stringer("account", address),
	)

	ch, err := s.auths.Challenge(*address)
	if err != nil {
		return "", err
	}

	if err := s.challenges.Put(*address, ch); err != nil {
		return "", err
	}
	return ch, nil
}

func (s *Service) Authorise(ctx context.Context, address *iam.Address, signature string) (string, error) {
	if address == nil {
		return "", errors.Bad(&errors.FieldViolation{
			Field:       "address",
			Description: "Missing account address",
		})
	}
	if signature == "" {
		return "", errors.Bad(&errors.FieldViolation{
			Field:       "signature",
			Description: "Missing signature",
		})
	}

	log.Trace(ctx, "auth.authorise", "Authorise",
		log.Stringer("account", address),
		log.String("signature", signature),
	)

	// Load current challenge for public key
	challenge, err := s.challenges.Get(*address)
	switch {
	case err == nil:
		// Good
	case errors.IsNotFound(err):
		return "", errors.PermissionDenied
	default:
		return "", err
	}

	// Verify challenge
	if err := s.auths.Verify(*address, challenge, signature); err != nil {
		return "", err
	}

	// Ensure account exists
	acc, err := s.accounts.Get(ctx, *address)
	if err != nil {
		return "", err
	}

	// Create and sign JWT
	token := jwt.NewWithClaims(jwtutil.StandardMethod, jwt.StandardClaims{
		Id:        acc.ID(),
		ExpiresAt: int64(utc.Now().Add(sessionTTL)),
	})
	signedToken, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
