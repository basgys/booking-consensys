package ethereum

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/basgys/booking-consensys/app/iam"
	"github.com/deixis/errors"
	"github.com/deixis/pkg/utc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Auth struct{}

func (a *Auth) Challenge(account iam.Address) (string, error) {
	b, err := genRandBytes(8)
	if err != nil {
		return "", errors.WithUnavailable(err, 0)
	}
	nonce := base64.StdEncoding.EncodeToString(b)
	return strings.Join([]string{"login", utc.Now().String(), nonce}, ":"), nil
}

func (a *Auth) Verify(account iam.Address, challenge string, hexSignature string) error {
	signature := common.FromHex(hexSignature)
	if len(signature) != crypto.SignatureLength {
		return fmt.Errorf(
			"invalid signature. It must be a hexadecimal representation of %d bytes, but got %d",
			crypto.SignatureLength,
			len(signature),
		)
	}

	// Recover public key that created the given signature
	sigPublicKey, err := crypto.SigToPub(signHash([]byte(challenge)), signature)
	if err != nil {
		return errors.WithPermissionDenied(err)
	}

	sigAccount := crypto.PubkeyToAddress(*sigPublicKey)
	if bytes.Equal(account.ETH().Bytes(), sigAccount.Bytes()) {
		return nil
	}
	return errors.PermissionDenied
}

// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L404
// signHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calculated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func genRandBytes(l int) ([]byte, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return nil, errors.Wrap(err, "rand error")
	}
	return b, nil
}
