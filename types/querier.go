package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query endpoints supported by the poa querier
const (
	QueryValidators = "validators"
	QueryValidator  = "validator"
	QueryParams     = "params"
)

// Defines the params for the following queries:
// - 'custom/poa/validator'
type QueryValidatorParams struct {
	ValidatorAddr sdk.ValAddress
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

// /*
// Below you will be able how to set your own queries:

// // QueryResList Queries Result Payload for a query
// type QueryResList []string

// // implement fmt.Stringer
// func (n QueryResList) String() string {
// 	return strings.Join(n[:], "\n")
// }

// */
