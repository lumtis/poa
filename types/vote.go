package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Structure to track the vote for:
// - An application to become validator
// - A proposal to kick a validator
type Vote struct {
	Subject   Validator        `json:"subject"`
	Approvals uint64           `json:"approvals"`
	Total     uint64           `json:"totals"`
	Voters    []sdk.AccAddress `json:"voter"`
}

// The subject of the vote
func (v Vote) GetSubject() Validator {
	return v.Subject
}

// The total number of approvals so far
func (v Vote) GetApprovals() uint64 {
	return v.Approvals
}

// THe total number of votes
func (v Vote) GetTotal() uint64 {
	return v.Total
}
