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
	Page, Limit   int
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress, page, limit int) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
		Page:          page,
		Limit:         limit,
	}
}

// Defines the params for the following queries:
// - 'custom/poa/validators'
type QueryValidatorsParams struct {
	Page, Limit int
}

func NewQueryValidatorsParams(page, limit int) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit}
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
