package ethereum_test

import (
	"fmt"
	"testing"

	"github.com/basgys/booking-consensys/app/auth/ethereum"
	"github.com/basgys/booking-consensys/app/iam"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestChallengeVerification(t *testing.T) {
	auth := ethereum.Auth{}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal("expect to generate key, but got", err)
	}
	addr := iam.Address(crypto.PubkeyToAddress(privKey.PublicKey))

	challenge, err := auth.Challenge(addr)
	if err != nil {
		t.Fatal("expect to get auth challenge, but got", err)
	}

	sign, err := crypto.Sign(signHash([]byte(challenge)), privKey)
	if err != nil {
		t.Fatal("expect to sign challenge, but got", err)
	}

	if err := auth.Verify(addr, challenge, common.Bytes2Hex(sign)); err != nil {
		t.Fatal("expect verify to succeed, but got", err)
	}
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
