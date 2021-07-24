package iam

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

type User struct {
	ID      string `json:"id"`
	GroupID string `json:"groupId"`
}

type Account struct {
	Address Address `json:"address"`
	UserID  string  `json:"userId"`
}

func (a *Account) ID() string {
	return a.Address.String()
}

type Group struct {
	ID  string `json:"id"`
	Ref string `json:"ref"`
}

type AccountProvider string

func (p AccountProvider) String() string {
	return string(p)
}

const (
	ProviderEthereum AccountProvider = "ethereum"
)

type Address common.Address

// ParseAddress parses an EIP55 hex string representation of an address
func ParseAddress(s string) (Address, error) {
	if !common.IsHexAddress(s) {
		return Address{}, errors.New("address is not a valid EIP55 hex string representation")
	}
	return Address(common.HexToAddress(s)), nil
}

func (a Address) ETH() common.Address {
	return common.Address(a)
}

// String returns an EIP55-compliant hex string representation of the address.
func (a Address) String() string {
	return common.Address(a).Hex()
}

// MarshalText implements TextMarshaler
func (a Address) MarshalText() (text []byte, err error) {
	return []byte(common.Address(a).Hex()), nil
}

// UnmarshalText implements TextUnmarshaler
func (a *Address) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	p, err := ParseAddress(string(text))
	if err != nil {
		return err
	}
	*a = p
	return nil
}
