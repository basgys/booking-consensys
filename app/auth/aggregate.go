package auth

import "github.com/basgys/booking-consensys/app/iam"

type Auth interface {
	// Challenge the given account `address`
	Challenge(address iam.Address) (string, error)
	// Verify checks whether the account that signed `challenge` with `signature`
	// matches the account `address`.
	//
	// When this function returns no errors, the authentication is considered
	// successful.
	Verify(address iam.Address, challenge string, signature string) error
}
