package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "poa"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

var (
	// Prefix for each key to a validator
	ValidatorsKey = []byte{0x21}

	// Prefix for each key to a validator index, by pubkey
	ValidatorsByConsAddrKey = []byte{0x22}
)

// Get the key for the validator with address
func GetValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorsKey, operatorAddr.Bytes()...)
}

// Get the key for the validator with pubkey
func GetValidatorByConsAddrKey(addr sdk.ConsAddress) []byte {
	return append(ValidatorsByConsAddrKey, addr.Bytes()...)
}
